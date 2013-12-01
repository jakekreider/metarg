METARG
======

*Half-baked METAR utility based on the [Debian metar package](http://packages.debian.org/wheezy/metar) by kees-guest, using to teach myself Go.*

Intended usage:  
`metarg KORD`  
`metarg -d KORD`  
*KORD being the airport code for Chicago O'Hare, where the weather always sucks*  

Will add more features as time and enthusiasm dictate.

TODO
----
  
*  Verbosity (`-v` flag)
*  Better parsing of units
*  Complete parsing of remarks
*  Better parsing of conditions
*  Decode METARs from input
*  Refactor this mess