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
Rodneys-iMac:gomap fogcat5$ appcfg.py -A commuteinfo-148920 -V v1 update ./
08:32 AM Application: commuteinfo-148920 (was: None); version: v1 (was: None)
08:32 AM Host: appengine.google.com
08:32 AM Starting update of app: commuteinfo-148920, version: v1
08:32 AM Getting current resource limits.
Your browser has been opened to visit:


    https://accounts.google.com/o/oauth2/auth?scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fappengine.admin+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcloud-platform+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2F&response_type=code&client_id=550516889912.apps.googleusercontent.com&access_type=offline


If your browser is on a different machine then exit and re-run this
application with the command-line parameter


  --noauth_local_webserver


Authentication successful.
08:32 AM Scanning files on local disk.
08:32 AM Cloning 69 application files.
08:32 AM Compilation starting.
08:32 AM Compilation: 63 files left.
08:33 AM Compilation completed.
08:33 AM Starting deployment.
08:33 AM Checking if deployment succeeded.
08:33 AM Deployment successful.
08:33 AM Checking if updated app version is serving.
08:33 AM Completed update of app: commuteinfo-148920, version: v1
</pre>

## Stuff For Routes
<pre>
  ../map/collector.py:@app.route('/collectdata')
  ../map/main.py:@app.route('/')
  ../map/main.py:@app.route('/sample')
  ../map/timeplot.py:@app.route('/drawday/<date>')
  ../map/timeplot.py:@app.route('/plotdatarev')
  ../map/timeplot.py:@app.route('/plotdatarev/<date>')
  ../map/timeplot.py:@app.route('/plotdata')
  ../map/timeplot.py:@app.route('/plotdata/<date>')
  ../map/timeplot.py:@app.route('/pmplot')
  ../map/timeplot.py:@app.route('/plot')
  ../map/timeplot.py:@app.route('/travel')
  ../map/timeplot.py:@app.route('/traveldata')
  ../map/timeplot.py:@app.route('/traveldata/<date>')
  ../map/timeplot.py:@app.route('/arrive')
  ../map/timeplot.py:@app.route('/arrivedata')
  ../map/timeplot.py:@app.route('/arrivedata/<date>')
  ../map/timeplot.py:@app.route('/settings', methods=['GET', 'POST'])
  ../map/timeplot.py:@app.route('/layout')
  ../map/timeplot.py:@app.route('/map')
  ../map/timeplot.py:@app.route('/')
  ../map/timeplot.py:@app.route('/whentogo')
</pre>


