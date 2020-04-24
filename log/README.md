# Log Log Log Log Log Log

I used to chose _golang/glog_ to log messages. 

It's a very nice package, but that's not what I want.

So I try to develop a tool to log messages in a way I want.

---

## Flags

- Use **-info** to turn on _log.Info*()_ prints

- Use **-warn** to turn on _log.Warning*()_ prints

- Use **-err false** to turn off _log.Error*()_ prints

- Use **-time** to print message with a time.Now().UTC().Format(time.RFC3339)[:19] tag.

### **Attention**

- _log.L*()_ and _log.Error*()_ prints are on by default.

---
## Print in info, warning and error styles:

### Use _Infof()_, _Warningf()_, _Errorf()_ to print informations **INLINE**, and use them just like you use _fmt.Printf()_.

#### Example[1]:

        1  package main
        2  import "github.com/achillesss/go-utils/log"
        3  func print(){
        4    log.Parse()
        5    log.Infof("My name is %s.", "Alex") // simply an information
        6    log.Warningf("Tell me how to spell %q, please.", "golang") // pretend to be warned
        7    log.Errorf("g-o-l-a-n.") // that's exactly an error
        8  }
        9  func main() {
        10     log.Parse()
        11     print()
        12 }

and you type 'go run main.go -info, -warn' in cmd-line. Then you'll get logs inline as below:

        [I_main.go_5] My name is Alex.    [W_main.go_6] Tell me how to spell "golang", please.    [E_main.go_7] g-o-l-a-n.


#### But why do this?

Just try tail -f _**your_log_name**.log_ | grep xxx in your cmd-line.

### Use _Infofln(), _Warningfln()_, _Errorfln()_ to print informations like _fmt.Printf(format+"\n", args...)_

Replace each _log.Infof_, _log.Warningf_, _log.Errorf_ by _log.Infofln_, _log.Warninfln_, _log.Errorfln_, and you'll get logs as below:

        [I_main.go_5] My name is Alex.
        [W_main.go_6] Tell me how to spell "golang", please.
        [E_main.go_7] g-o-l-a-n.

#### But why do this?

Want you to debug with all messages **inline**? I don't know and I don't want.

### Use **_-time_** to print with a time tag by _time.Now()_ in formation of "2016-01-02T03:04:05"

if you turn on it, your all logs will turn into the form as below:
        
        [X_xxx.go_line 2016-01-02T03:04:05] xxxxxxx
---
## Log:

### Use _L()_, to log a function response

#### Example[2]:
        1  package main
        2  import "github.com/achillesss/go-utils/log"
        3  type response struct {}
        4  func (r *response) String() (res string) {
        5     // marshal r into res
        6     return 
        7  }
        8  func getResponse() (res *response, err error) {
        9      log.Infof("start to get response.")
        10     // get response
        11     log.L(err != nil, "%s failed. Error: %v. Resp: %s",log.FuncName(), err, res.String*())
        13 }
        14 func main() {
        15     log.Parse()
        16     getResponse()            
        17 }

type 'go run main.go' and

if you've just got a bad response with an error, you'll get a log as below:

    [I_main.go_9] start to get response.    [L_main.go_11] getResponse failed. Error: xxx. Resp: xxx    [I_main.go_12] finish get response.

and if you replace _log.L()_ with _log.Lln()_, _log.Infof()_ with _log.Infofln()_ following logs is what you'll get:
    
    [I_main.go_9] start to get response.
    [L_main.go_11] getResponse failed. Error: xxx. Resp: xxx
    [I_main.go_12] finish get response.

---
## Error Formating:

Time flies. I'll fill up this field later.

---
## Print name of your _func()_:

Time flies. I'll fill up this field later.

---
## Versions:
* #### 1.0:
    Add basic funtions
* #### 2.0:
    Rebuild and add some toggles
* #### 2.0.1:
    Add README.md
* #### 2.1.0
    Add flag '-time' to turn on a time.Now().UTC() tag formatted as "_2006-01-02T03:04:05_" 
