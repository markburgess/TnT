//
// Copyright © Mark Burgess
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// ***************************************************************************
//*
//* TnT - Trust and Trustworthiness
//* Applied to Trust and Semantics Learning
//*
// ***************************************************************************

package TnT

import (
	"strings"
	"fmt"
	"regexp"
	"os"
	"time"
	"sort"
	"unicode"
	"math"
)

// **********************************************************************

const NANO = 1000000000
const MILLI = 1000000
const NOT_EXIST = 0

const KVDIR = "/tmp/TnT_KV/"

// ***************************************************************************

type PromiseHistory struct {

	// Use this as an event tracker, CFEngine style

	PromiseId string     `json:"_key"`

	// Three points for derivative

	Q         float64    `json:"q"`
	Q1        float64    `json:"q1"`
	Q2        float64    `json:"q2"`

	Q_av      float64    `json:"q_av"`
	Q_var     float64    `json:"q_var"`

	T         int64      `json:"lastT"`
	T1        int64      `json:"lastT1"`
	T2        int64      `json:"lastT2"`

	Dt_av     float64    `json:"dT"`
        Dt_var    float64    `json:"dT_var"`

	V         float64    `json:"V"`
	AntiT     float64    `json:"antiT"`

	Units     string     `json:"units"`
}

// *********************************************************************

type PromiseContext struct {

	Time  time.Time
	Name  string
	Plock Lock
}

// **********************************************************************
// Promise Context
// **********************************************************************

func PromiseContext_Begin(name string) PromiseContext {

	before := time.Now()
	return StampedPromiseContext_Begin(name, before)
}

// **********************************************************************

func StampedPromiseContext_Begin(name string, before time.Time) PromiseContext {

	// Set up memory for history, register callbacks

	var ctx PromiseContext
	ctx.Time = before
	ctx.Name = KeyName(name,0)

	// *** begin ANTI-SPAM/DOS PROTECTION ***********

	ifelapsed := int64(30)   // these params should be policy
	expireafter := int64(60)

	now := time.Now().UnixNano()

	ctx.Plock = BeginService(name,ifelapsed,expireafter, now) 

	// *** end ANTI-SPAM/DOS PROTECTION ***********

	return ctx
}

// **********************************************************************

func PromiseContext_End(ctx PromiseContext) PromiseHistory {

	after := time.Now()
	return StampedPromiseContext_End(ctx,after)
}

// **********************************************************************

func StampedPromiseContext_End(ctx PromiseContext, after time.Time) PromiseHistory {

	before := ctx.Time

	EndService(ctx.Plock)

	const collname = "conn"
	var key string

	// Semantic donut time key ..

	_, timeslot := DoughNowt(time.Now())
	
	if ctx.Name == "" {
		key = timeslot
	} else {
		key = ctx.Name+":"+timeslot
	}
	
	// make b = promise execution interval (latency) in this case

	b := float64(after.Sub(before)) // time difference now-previous

	// Direct db writes, these are separated from the time-based averaging

	previous_value := GetKV(collname,ctx.Name+"latency")
	previous_time := GetKV(collname,ctx.Name+"lasteen")

	var dt,db float64

	if previous_time.V == 0 {
		dt = 300 * NANO  // default bootstrap
	} else {
		dt = float64(after.UnixNano()) - previous_time.V
	}

	if previous_value.V == 0 {
		db = b/2         // default bootstrap
	} else {
 		db = b - previous_value.V
	}

	dtau := dt/db * b

	e := LearnUpdateKeyValue(collname,key,time.Now().UnixNano(),b,"ns")

	var lastlatency,lasttime KeyValue

	// Make the values latency

	lastlatency.K = ctx.Name+"latency"
	lastlatency.V = b

	lasttime.K = ctx.Name+"lastseen"
	lasttime.V = float64(after.UnixNano())

	fmt.Println("------- INSTRUMENTATION --------------")

	AddKV(collname,lastlatency)
	AddKV(collname,lasttime)

	fmt.Println("   Location:", ctx.Name+collname)
	fmt.Println("   Promise duration b (ms)", e.Q/MILLI,"=",b/MILLI)
	fmt.Println("   Running average 50/50", e.Q_av/NANO)

	fmt.Println("   Change in promise since last sample",db)
	fmt.Println("   Promise derivative b/s", db/dt)
	fmt.Println("")
	fmt.Println("   Time since last sample (s) phase",dt/NANO)
	fmt.Println("   Time signal uncertainty dtau (s) group",dtau/NANO)
	fmt.Println("   Running average sampling interval",e.Dt_av/NANO)
	fmt.Println("------- INSTRUMENTATION --------------")
	return e
}

// **********************************************************************

func AssessPromiseOutcome(e PromiseHistory, assessed_quality,promise_upper_bound,trust_interval float64) float64 {

	promised_ns := promise_upper_bound * NANO
	trust_ns := trust_interval * NANO

	key := e.PromiseId

	// This function decides the kinetic trust and adjusts the potential
	// V based on real time promise keeping. It doesn't consider the initial
	// determination of V -- i.e. whether we want to talk to the other agent
	// at all (as in security)

	var sig float64 = math.Sqrt(e.Q_var)

	// Here we've measured the timing and we've looked at the content
	// Now we need to determine the promise-kept assessment degree
	// and adjust the long term history for this promise

	// The trouble is that we don't usually know what was promised...

	promise_level := 1/(1+math.Exp(3*(e.Q-promised_ns)/promised_ns))

	fmt.Println("Promise level",promise_level,"+-",sig/promised_ns,"raw",e.Q/NANO,promise_upper_bound)

	if e.Dt_av == 0 {
		e.Dt_av = 1.0
	}

	fmt.Println("Assessing expected sampling interval",float64(e.T)/e.Dt_av)
	fmt.Println("Assessing desired sampling interval",float64(e.T)/trust_ns)

	// The assessed payload is the user defined arbitrary up or downvote
	// How well did we keep our promise payload?

	fmt.Println("Assessing expected Q level",float64(e.Q)/e.Q_av)
	fmt.Println("Assessing desired Q level",float64(e.Q)/promised_ns)
	fmt.Println("Assessing payload",assessed_quality)
	fmt.Println("Assessing level change",(e.Q-e.Q1)/promised_ns)

	// Get our previous estimate of reliability

	reliability := GetKV("PromiseKeeping",key)

	if reliability.V == 0 {

		reliability.V = 0.5 // Start evens
	}

	// Q is always positive (latency here...)
 	// Some assessments of the event's general timeliness
	// A significant timescale for latency is 0.1 second?

	delta := promise_level * assessed_quality

	if math.Abs(e.Q_av) < sig {  // Down vote for noisy behaviour

		fmt.Println("1.PENALTY!")
		delta = delta / 1.5
	}

	// derivatives are possible signs of stress / coping (confidence)
	// if first first second derivatives are growing, this is not good for latency

	dqdt := FirstDerivative(e,promised_ns,trust_ns)
	d2qdt2 := SecondDerivative(e,promised_ns,trust_ns)

	const sensitivity = 0.01 // should this be the same for 1st and second?

	if dqdt < -sensitivity {
		fmt.Println("Gradient reducing (spot measure)")
		delta = delta + 0.1
	} else if dqdt > sensitivity {
		fmt.Println("Gradient increasing (spot measure)")
		delta = delta - 0.1
		fmt.Println("2.PENALTY!")
	}

	if d2qdt2 < -sensitivity {
		fmt.Println("Curvature decelerating (positive force)")
		delta = delta + 0.1
	} else if d2qdt2 > sensitivity {
		fmt.Println("Curvature accelerating (negative force)")
		delta = delta - 0.1
		fmt.Println("3.PENALTY!")
	}

	//if math.Fabs(SecondDeriv(e)) > SCALE {
	//	delta = delta - 0.5
	//}

	// Adjust reliability according to timing AND quality

	fmt.Println("Old ML running reliability(delta)",reliability.V)

	if delta < 0 {

		delta = 0
	}

	reliability.K = key
	reliability.V = reliability.V * 0.4 + delta * 0.6

	fmt.Println("New ML running reliability(delta)",reliability.V,delta)

	AddKV("PromiseKeeping",reliability)

	return reliability.V
}

// ****************************************************************************
// Heuristic context, CFEngine style
// ****************************************************************************

var CONTEXT map[string]float64

// *******************************************************************************

func ContextActive(s string) {

	// Machine learn in a Bayesian fashion a context state assumed true if called

	CONTEXT[s] = 0.5 + 0.5 * CONTEXT[s]
}

// *******************************************************************************

func ContextSet() []string {

	var result []string

	for s := range CONTEXT {
		if CONTEXT[s] > 0 {
			result = append(result,s)
		}
	}

	sort.Strings(result)

	return result
}

// *******************************************************************************

func InitializeContext() {

	// Reset / empty all signal values in context

	if !IsDir(KVDIR) {
		
		os.MkdirAll(KVDIR, 0755)
	}

	CONTEXT = make(map[string]float64)
}

// *******************************************************************************

func SetContext(s string,c float64) {

	// Set the probability / confidence of the identifer explicitly

	CONTEXT[s] = c
}

// *******************************************************************************

func IsDefinedContext(s string) bool {

	// Evalute general boolean expressions CFEngine style

	r,confidence := ContextEval(s)

	if r == "bad expression" {
		fmt.Println("Bad context expression:",s)
	}

	return (confidence > 0)
}

// *******************************************************************************

func Confidence(s string) float64 {

	// Evalute a real number return for expressions CFEngine style

	_,confidence := ContextEval(s)

	return confidence
}

// ***********************************************************************

func ContextEval(s string) (string,float64) {

	// Return an estimated confidence in the quasi-Boolean expression s

	expr := CleanExpression(s)

	or_parts := SplitWithParensIntact(expr,'|')

	or_result := 0.0

	for or_frag := range or_parts {

		and_parts := SplitWithParensIntact(or_parts[or_frag],'.') 

		and_result := 1.0

		for and_frag := range and_parts {

			if s[0] == '(' && and_parts[and_frag] == s {
				fmt.Println("\nIrreducible context expression: ",s,"\n")
				return "bad expression", -1.0
			}

			var res float64
			token := strings.TrimSpace(and_parts[and_frag])
			
			switch token[0] {

			case '!':
				switch token[1] {
				case '(': 
					_,res = ContextEval(token[1:])
				default:
					res = CONTEXT[token[1:]]
				}

				if res > 0 {
					res = 0
				} else {
					res = 1
				}

			case '(': 
				_,res = ContextEval(token)
			default:
				res = CONTEXT[token]
			}

			// P(A and B) = AB
			and_result *= res
		}

		// P(A or B) ~ (A + B - AB)
		or_result += and_result - (or_result) * and_result
	}

	return expr,or_result
}

// ***********************************************************************

func CleanExpression(s string) string {

	s = TrimParen(s)
	r1 := regexp.MustCompile("[|]+") 
	s = r1.ReplaceAllString(s,"|") 
	r2 := regexp.MustCompile("[&]+") 
	s = r2.ReplaceAllString(s,".") 
	r3 := regexp.MustCompile("[.]+") 
	s = r3.ReplaceAllString(s,".") 

	return s
}

// ***********************************************************************

func SplitWithParensIntact(expr string,split_ch byte) []string {

	var token string = ""
	var set []string

	for c := 0; c < len(expr); c++ {

		switch expr[c] {

		case split_ch:
			set = append(set,token)
			token = ""

		case '(':
			subtoken,offset := Paren(expr,c)
			token += subtoken
			c = offset-1

		default:
			token += string(expr[c])
		}
	}

	if len(token) > 0 {
		set = append(set,token)
	}

	return set
} 

// ***********************************************************************

func Paren(s string, offset int) (string,int) {

	var level int = 0

	for c := offset; c < len(s); c++ {

		if s[c] == '(' {
			level++
			continue
		}

		if s[c] == ')' {
			level--
			if level == 0 {
				token := s[offset:c+1]
				return token, c+1
			}
		}
	}

	return "bad expression", -1
}


// ***********************************************************************

func TrimParen(s string) string {

	var level int = 0
	var trim = true

	if len(s) == 0 {
		return s
	}

	s = strings.TrimSpace(s)

	if s[0] != '(' {
		return s
	}

	for c := 0; c < len(s); c++ {

		if s[c] == '(' {
			level++
			continue
		}

		if level == 0 && c < len(s)-1 {
			trim = false
		}
		
		if s[c] == ')' {
			level--

			if level == 0 && c == len(s)-1 {
				
				var token string
				
				if trim {
					token = s[1:len(s)-1]
				} else {
					token = s
				}
				return token
			}
		}
	}
	
	return s
}

// ***************************************************************************
// Key-Value storage
// ***************************************************************************

type KeyValue struct {

	K  string  `json:"_key"`
	V  float64 `json:"value"`
}

// ***************************************************************************

func AddKV(collname string,kv KeyValue) {

	if !IsDir(KVDIR) {
		
		os.MkdirAll(KVDIR, 0755)
	}

	s := fmt.Sprintf("%+v",kv)

	data := []byte(s)
	err := os.WriteFile(KVDIR+collname+kv.K, data, 0644)

	if err != nil {
		fmt.Println("Unable to write promise",kv,err)
		os.Exit(-1)
	}

}

// **************************************************

func GetKV(collname,key string) KeyValue {

	var kv KeyValue

	b,_ := os.ReadFile(KVDIR+key)

	fmt.Sscanf(string(b),"%f",&kv.V)
	kv.K = key
	kv.V = CONTEXT[key]
	return kv
}

// **************************************************
// Promise tracking
// **************************************************

func AddPromiseHistory(collname, key string, e PromiseHistory) {

	if !IsDir(KVDIR) {
		
		os.MkdirAll(KVDIR, 0755)
	}

	s := fmt.Sprintf("%+v",e)

	data := []byte(s)
	err := os.WriteFile(KVDIR+collname+key, data, 0644)

	if err != nil {
		fmt.Println("Unable to write promise",key,err)
		os.Exit(-1)
	}
}

// **************************************************

func GetPromiseHistory(collname, key string) (bool,PromiseHistory) {

	var v PromiseHistory

	data,err := os.ReadFile(KVDIR+collname+key)

	fmt.Sscanf(string(data),"%+v",&v)

	if err != nil {
		return true, v
		
	} else {
		var dud PromiseHistory
		dud.T = NOT_EXIST
		dud.Q = NOT_EXIST
		return false, dud		
	}
}

// **************************************************

func LearnUpdateKeyValue(collname,key string, now int64, q float64, units string) PromiseHistory {

	// now should be time.Now().UnixNano()

	var e PromiseHistory

	e.PromiseId = key

	// Slide derivative window

	// time is weird in go. Duration is basically int64 in nanoseconds

	exists, previous := GetPromiseHistory(collname,key)
	
	if !exists {

		// Initial bootstrap defaults

		e.Q_av = 0.6 * float64(q)
		e.Q_var = 0

		e.T = now
		e.Dt_av = 0
		e.Dt_var = 0

		AddPromiseHistory(collname,key,e)

	} else {
		e.Q2 = previous.Q1
		e.Q1 = previous.Q
		e.Q = q

		e.Units = units

		e.Q_av = 0.5 * previous.Q + 0.5 * float64(q)
		dv2 := (e.Q-e.Q_av) * (e.Q-e.Q_av)
		e.Q_var = 0.5 * e.Q_var + 0.5 * dv2
		
		e.T2 = previous.T1
		e.T1 = previous.T
		e.T = now

		dt := float64(now-previous.T) // time difference now-previous

		e.Dt_av = 0.5 * previous.Dt_av + 0.5 * dt
		e.Dt_var = 0.5 * e.Q_var + 0.5 * (e.Dt_av-dt) * (e.Dt_av-dt)

		AddPromiseHistory(collname,key,e)
	}

	return e
}

// ****************************************************************************

func FirstDerivative(e PromiseHistory, qscale,tscale float64) float64 {

	dq := (e.Q - e.Q1)/qscale
	dt := float64(e.T-e.T1)/tscale

	if dt == 0 {
		return 0
	}

	dqdt := dq/dt

	fmt.Println("Deriv dq/dt (latency)",dqdt)

	return dqdt
}

// ****************************************************************************

func SecondDerivative(e PromiseHistory, qscale,tscale float64) float64 {

	dv := ((e.Q - e.Q1)/float64(e.T-e.T1) - (e.Q1 - e.Q2)/float64(e.T1-e.T2))/qscale*tscale

	dt := (e.Q1 *float64(e.T-e.T1)/tscale)

	d2qdt2 := dv/dt

	if dt == 0 {
		return 0
	}

	fmt.Println("Deriv d2q/dt2 (latency)",d2qdt2)

	return d2qdt2
}

// ****************************************************************************
// Semantic 2D time
// ****************************************************************************

var GR_DAY_TEXT = []string{
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
        "Sunday",
    }
        
var GR_MONTH_TEXT = []string{
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
}
        
var GR_SHIFT_TEXT = []string{
        "Night",
        "Morning",
        "Afternoon",
        "Evening",
    }

// For second resolution Unix time

const CF_MONDAY_MORNING = 345200
const CF_MEASURE_INTERVAL = 5*60
const CF_SHIFT_INTERVAL = 6*3600

const MINUTES_PER_HOUR = 60
const SECONDS_PER_MINUTE = 60
const SECONDS_PER_HOUR = (60 * SECONDS_PER_MINUTE)
const SECONDS_PER_DAY = (24 * SECONDS_PER_HOUR)
const SECONDS_PER_WEEK = (7 * SECONDS_PER_DAY)
const SECONDS_PER_YEAR = (365 * SECONDS_PER_DAY)
const HOURS_PER_SHIFT = 6
const SECONDS_PER_SHIFT = (HOURS_PER_SHIFT * SECONDS_PER_HOUR)
const SHIFTS_PER_DAY = 4
const SHIFTS_PER_WEEK = (4*7)

// ****************************************************************************
// Semantic spacetime timeslots
// ****************************************************************************

func DoughNowt(then time.Time) (string,string) {

	// Time on the torus (donut/doughnut) (CFEngine style)
	// The argument is a Golang time unit e.g. then := time.Now()
	// Return a db-suitable keyname reflecting the coarse-grained SST time
	// The function also returns a printable summary of the time

	year := fmt.Sprintf("Yr%d",then.Year())
	month := GR_MONTH_TEXT[int(then.Month())-1]
	day := then.Day()
	hour := fmt.Sprintf("Hr%02d",then.Hour())
	mins := fmt.Sprintf("Min%02d",then.Minute())
	quarter := fmt.Sprintf("Q%d",then.Minute()/15 + 1)
	shift :=  fmt.Sprintf("%s",GR_SHIFT_TEXT[then.Hour()/6])

	//secs := then.Second()
	//nano := then.Nanosecond()

	dayname := then.Weekday()
	dow := fmt.Sprintf("%.3s",dayname)
	daynum := fmt.Sprintf("Day%d",day)

	// 5 minute resolution capture
        interval_start := (then.Minute() / 5) * 5
        interval_end := (interval_start + 5) % 60
        minD := fmt.Sprintf("Min%02d_%02d",interval_start,interval_end)

	var when string = fmt.Sprintf("%s,%s,%s,%s,%s at %s %s %s %s",shift,dayname,daynum,month,year,hour,mins,quarter,minD)
	var key string = fmt.Sprintf("%s:%s:%s",dow,hour,minD)

	return when, key
}

// ****************************************************************************

func GetUnixTimeKey(now int64) string {

	// Time on the torus (donut/doughnut) (CFEngine style)
	// The argument is in traditional UNIX "time_t" unit e.g. then := time.Unix()
	// This is a simple wrapper to DoughNowt() returning only a db-suitable keyname

	t := time.Unix(now, 0)
	_,slot := DoughNowt(t)

	return slot
}

// ****************************************************************************

func GetAllWeekMemory(collname string) []float64 {

	// Used in Machine Learning of weekly patterns, keys labelled with DoughNowt()
	// Returns a vector from Monday morning 00:00 to Sunday evening 11:55 in 5 min grains
	// The collection name is assumed to point to an Arango KeyValue database collection

	var now int64
	var data []float64

	for now = CF_MONDAY_MORNING; now < CF_MONDAY_MORNING + SECONDS_PER_WEEK; now += CF_MEASURE_INTERVAL {

		slot := GetUnixTimeKey(now)
		kv := GetKV(collname,slot)
		data = append(data,kv.V)
	}

	return data
}

// ****************************************************************************

func SumWeeklyKV(collname string,t int64, value float64){

	// Create a cumuluative weekly periodogram database KeyValue store
	// the time t should be in time.Unix() second resolution

	key := GetUnixTimeKey(t)
	kv := GetKV(collname,key)
	kv.K = key
	kv.V = value + kv.V
	AddKV(collname,kv)
}

// ****************************************************************************

func LearnWeeklyKV(collname string,t int64, value float64){

	// Create an averaging weekly periodogram database KeyValue store
	// the time t should be in time.Unix() second resolution

	key := GetUnixTimeKey(t)
	kv := GetKV(collname,key)
	kv.K = key
	kv.V = 0.5 * value + 0.5 * kv.V
	AddKV(collname,kv)
}

// ****************************************************************************

func AddWeeklyKV_Unix(collname string, t int64, value float64) {

	// Add a single key value to a weekly periodogram, update by Unix() time key

	var kv KeyValue
	kv.K = GetUnixTimeKey(t)
	kv.V = value
	AddKV(collname,kv)
}

// ****************************************************************************

func AddWeeklyKV_Go(collname string, t time.Time, value float64) {

	// Add a single key value to a weekly periodogram, update by Golang time.Time key

	var kv KeyValue
	_,kv.K = DoughNowt(t)
	kv.V = value
	AddKV(collname,kv)
}

// *****************************************************************
// Adaptive locks ...
// Pedagogical implementation of self-healing locks as used in CFEngine
// We need a 1:1 unique name for client requests and lock names
// Also, it's important to write service code that's interruptible, especially
// in golang where you can't forcibly signal by imposition as with preemptive MT
//
//  lock := BeginService(...)
//    ...
//  EndService(lock)
// *****************************************************************

const LOCKDIR = "/tmp/TnT_Locks/" // this should REALLY be a private location
const NEVER = 0

type Lock struct {

	Ready bool
	This  string
	Last  string
}

// *****************************************************************

func BeginService(name string, ifelapsed,expireafter int64, now int64) Lock {

	var lock Lock

	lock.Last = fmt.Sprintf("last.%s",name)
	lock.This = fmt.Sprintf("lock.%s",name)
	lock.Ready = true
	
	lastcompleted := GetLockTime(lock.Last)

	elapsedtime := (now - lastcompleted) / NANO // in seconds

	if (elapsedtime < ifelapsed) {

		fmt.Println("Too soon since last",lock.Last,elapsedtime,"/",ifelapsed)
		lock.Ready = false
		return lock
	}

	starttime := GetLockTime(lock.This)

	if (starttime == NEVER) {

		//Println("No running lock...")

	} else {

		runtime := (now-starttime) / NANO

		if (runtime > expireafter) {

			// server threads can't be forced to quit, 
			// so we can only ask nicely to release resources
			// as part of a standard promise
			// If the thread can change something downstream, it needs to be stopped
			// For a read only server process, it's safe to continue

			RemoveLock(lock.This)
		}
	}

	AcquireLock(lock.This)
	return lock
}

// *****************************************************************

func EndService(lock Lock) {

	RemoveLock(lock.This)
	RemoveLock(lock.Last)
	AcquireLock(lock.Last)
}

// *****************************************************************

func GetLockTime(filename string) int64 {

	fileinfo, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {

			return NEVER

		} else {
			fmt.Println("Insufficient permission",err)
			os.Exit(1)
		}
	}

	return fileinfo.ModTime().UnixNano()
}

// *****************************************************************

func AcquireLock(name string) {

	if !IsDir(LOCKDIR) {
		
		os.MkdirAll(LOCKDIR, 0755)
	}

	f, err := os.OpenFile(LOCKDIR+name,os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Couldn't acquire lock to create",name,err)
		return
	}

	f.Close()
}

// *****************************************************************

func RemoveLock(name string) {

	err := os.Remove(LOCKDIR+name)

	if err != nil {
		//fmt.Println("Unable to remove",name,err)
	}
}

//**************************************************************

func KeyName(s string,n int) string {

	strings.Trim(s,"\n ")
	
	if len(s) > 40 {
		s = s[:40]
	}

	var key string

	runes := []rune(s)

	for r := range runes {

		if !unicode.IsPrint(runes[r]) {
			runes[r] = 'x'
		}
	}

	m := regexp.MustCompile("[^a-zA-Z0-9]") 
	str := m.ReplaceAllString(string(runes),"-") 

	if n > 0 {
		key = fmt.Sprintf("%s_%d",str,n)
	} else {
		key = str
	}

	return strings.ToLower(key)
}

//**************************************************************

func CanonifyName(s string) string {

	return KeyName(s,0)
}

//**************************************************************
// File stat
//**************************************************************

func IsFile(filename string) bool {

	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

//**************************************************************

func IsDir(filename string) bool {

	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}