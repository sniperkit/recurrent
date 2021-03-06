# Recurrent

[![GoDoc](https://godoc.org/github.com/efritz/recurrent?status.svg)](https://godoc.org/github.com/efritz/recurrent)
[![Build Status](https://secure.travis-ci.org/efritz/recurrent.png)](http://travis-ci.org/efritz/recurrent)
[![Maintainability](https://api.codeclimate.com/v1/badges/0ad38bd1583c57173ebb/maintainability)](https://codeclimate.com/github/efritz/recurrent/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/0ad38bd1583c57173ebb/test_coverage)](https://codeclimate.com/github/efritz/recurrent/test_coverage)

Go library for periodically scheduling a task.

## Example

First, you must create a scheduler and pass it a reference to the function that should
execute periodically. It also takes an interval, which is the amount of time that must
pass before executing the function again. Schedulers begin once the `Start` method is
called. A scheduler should be stopped (explicitly or by `defer`) to prevent goroutine
leaks.

```go
sched := NewScheduler(cronTask, WithInterval(time.Second * 5))

sched.Start()
defer sched.Stop()
```

The function can be executed on-demand from within the scheduler by sending it a signal
to execute. This call is always non-blocking.

```go
sched.Signal()
```

A *throttled* scheduler can also be created, which throttles the frequency at which an
explicit signal to execute the function has an effect on the scheduler. The `cronTask`
will be executed at most once per-second, regardless of how many signals are received.

```go
sched := NewScheduler(cronTask, WithInterval(time.Second * 5), WithThrottle(time.Second))

sched.Start()
defer sched.Stop()

for {
    sched.Signal()
}
```

## License

Copyright (c) 2016 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
