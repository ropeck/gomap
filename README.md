# Gomap
  A simple project to use the Google maps API and show traffic commute times.
  Translated from the same project written in Python.

# (TODO)
## parts to convert
These items still need to be converted from the python version of mapapp.

# Overview

## Dependencies
  This module uses the directions module which does the actual API calls.
It's in github.com/ropeck/directions and included in the appengine app.

## Deployment

To deploy an appengine package, use this

<pre>
appcfg.py -A mapappgo -V v1 update ./
<pre>

## Testing
App engine SDK runs unit tests in the *_test.go files like this:
<pre>
goapp test
</pre>

## Stuff TODO

* cleanup text display in the output for times and distances

 

