## Goal

This demo demonstrate how to use pprof to diagnose CPU usage

The goal of these programs is to generate a HTTP request according to a Req struct and template data.

## Intro

The 'old' program maintains some templates for URL, body, each header and query. Those templates will be executed to generate settings of a HTTP request.

The 'new' program use only one template which is JSON encoded Req struct. It will be executed and then decoded into Req struct with actual settings.

Tool 'pprof' can be used to analyse the performance of these two methods.

## How to

Enter into each sub-directory and run command:

```
go test -benchmem -run=^$ -bench ^BenchmarkSequetiallyGenRequest$ -cpuprofile cpu.prof
```

And start 'pprof' tool like following:

```
go tool pprof old.test cpu.prof
```

Visualize graph through web browser:

```
(pprof) web
```

It depends on Graphviz. Install it if you got an warning.

It can be found that the 'new' method is slower than 'old' one. The call of json.Unmarshal is the bottleneck. 
