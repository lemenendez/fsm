# Finite State Machine Implementation (FSM)

## The Idea

A Finite State Machine has multiple use cases, the state of an entity for example a *User*. Another example is the state of a *Tenant* for a Saas company, the state of the executing of a long-running task, or the state of a money transfer.

## Definition

- The system must be describable by a finite set of states.
- The system must have a finite set of inputs and/or events that can trigger transitions between states.
- The behavior of the system at a given point in time depends upon the current state and the input or the event that occurs at that time.
- For each state the system may be in, the behavior is defined for each possible input or event.
- The system has a particular initial state.

## Naming Rules

1. The State name must be UPPERCASE
2. The State name must not contain spaces
3. The State name must not contain numbers
4. The Transition name must be UPPERCASE
5. The Transition name must not contain spaces
6. The Transition name must not contain numbers

## Examples

A complete set of examples are under the `test` folder.

### Simple Enable/Disable

The next code declares 2 states: ENABLED and DISABLED, and 2 transitions ENABLE and DISABLE

```GO
const ENABLED = "ENABLED"
const DISABLED = "DISABLED"
const ENABLE = "ENABLE"
const DISABLE = "DISABLE"
```

The next lines declare a finite state machine (fsm), then it adds its states and its transitions.

```GO
f := fsm.NewFSM("Basic Disabled/Enabled")
f.AddState(DISABLED)
f.AddState(ENABLED)
f.Init(DISABLED)

f.AddTrans(ENABLED, DISABLED, DISABLE)
f.AddTrans(DISABLED, ENABLED, ENABLE)
```

Now we declare the function `myFunc`, that funcion will be executed after fsm does a transition between the states.
Note: It is always a good idea to do this:

```GO
err := f.AddState(DISABLE)
if er...
```

```GO
myFunc := func(pre string, cur string, action string) {
    t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
}

err := f.Exec(ENABLE, ENABLED, myFunc)

if err != nil {
    t.Log(err)
    t.Errorf("should allow transicion")
}
t.Log(f.GetTrans())
```

### SaaS Plan

![A test image](docs/README.png)

A user starts in __Trial__, __Basic__ or __Premium__. Once in __Basic__ can *Upgrade* to __Premium__. Once in __Premium__ can *Downgrade* to __Basic__. In __Trial__ can *upgrade* to __Basic__ or __Premium__. For this example we're not going to use the __Expired__ state.

<!--
```
@startuml
[*] --> Trial
[*] --> Basic
[*] --> Premium
Trial --> [*]
Trial -> Basic
Trial -> Premium
Basic -> Premium
Premium -> Basic
Premium --> [*]
Basic --> [*]
@enduml
```
-->

Define a Plan struct:

```GO
type Plan struct {
    Id    int32
    Name  string
    State *fsm.FSM
}
```

Define the 'Constructor'

```GO
func NewPlan(plan string) (*Plan, error) {

    f := fsm.NewFSM("SAAS Account State V1.0")

    f.AddState("TRIAL")
    f.AddState("BASIC")
    f.AddState("PREMIUM")

    f.AddTrans("TRIAL", "BASIC", "UPGRATE")
    f.AddTrans("TRIAL", "PREMIUM", "UPGRATE")
    f.AddTrans("BASIC", "PREMIUM", "UPGRATE")
    f.AddTrans("PREMIUM", "BASIC", "DOWNGRATE")

    c := &Plan{
        Id:    0,
        Name:  "Standard",
        State: f,
    }
    err := f.Init(plan)
    if err == nil {
        return c, nil
    }
    return nil, err
}
```

We use the __Init__ function to initialize the fsm to its initial state.

Now we define the Upgrade and Downgrade helper functions.

```GO
func (a *Plan) Upgrade(name string) bool {
    err := a.State.Exec("UPGRATE", name, a.stateTransitionHanlder)
    if err != nil {
        a.t.Logf(err.Error())
    }
    return err == nil
}
```

```GO
func (a *Plan) Downgrate(name string) bool {
    err := a.State.Exec("DOWNGRATE", name, a.stateTransitionHanlder)
    if err != nil {
        a.t.Logf(err.Error())
    }
    return err == nil
}
```

The __stateTransitionHandler__  function will be called after the fsm validate the execution.

```GO
func (a *Plan) stateTransitionHandler(pre string, cur string, action string) {
    a.t.Logf("Previous State:%v, New State:%v, Action:%v", pre, cur, action)
}
```

Now finally we use it. We use *NewPlan* funcion to create and also to initilize our fsm.Then we call the *Upgrade* funcion to check if we can or cannot upgrade the plan.

```GO
    p, err := NewPlan("TRIAL")
    if p != nil {
        if p.Upgrade("BASIC") {
            t.Logf("New State:%v", a.State.GetState())
        } else {
            t.Errorf("cannot updagrade")
        }
    } else {
        t.Errorf(err.Error())
    }
```

## Docker

### Build

`docker build -t fsm .`

### Running

`docker run -it --rm -v "$PWD":/go/src fsm bash`

### Testing

`go test ./test/ -v`

Testing with coverage: `go test ./test/ -v -coverprofile=cover.out -coverpkg=.`

Testing with tool: `go tool cover -html=$PWD/cover.out -o $PWD/cover.html`

## Math Definition

A finite automaton M is defined by a 5-tuple (Σ, Q, q 0 , F, δ), where

- Σ is the set of symbols representing input to M
- Q is the set of states of M
- q 0 ∈ Q is the start state of M
- F ⊆ Q is the set of final states of M
- δ : Q × Σ → Q is the transition function

Refer to the [fsm doc](docs/fsm-notes.pdf)

