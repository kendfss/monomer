## monomer

A VST2 Plugin that enables the user to convert stereo signals to mono and/or
swap the left and right channels.

It's heavily indebted by [DFX](http://destroyfx.org/docs/monomaker.html)'s
Monomaker, in fact I originally wrote this because I needed something to replace
it when I wanted to finish some tracks on a mac that didn't support 32bit
plugins.

# compiling

to compile it for your current platform just use `go build --tags=plugin` and
you'll get a shared/dynamic object/library in your current directory.

to compile for windows from another platform use:

```bash
CGO_ENABLED=1 GOOS=windows CC=path/or/name/of/your/windows/c/compiler go build --buildmode=c-shared --tags=plugin -o monomaker.dll
```
