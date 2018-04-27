# syncmap - a Mutex-protected Map

I'm sure this already exists somewhere else, and in a better form...but...
I keep finding myself needing a concurrent Map. And rather than continually
copy-pasting the same code, I'm going to put it here.

This type combines a sync.RWMutex with a map. 
