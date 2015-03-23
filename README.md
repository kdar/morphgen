MorphGEN
=========

Download the binary distribution here: http://github.com/kdar/morphgen-binary

An application that helps in generating TMorph codes for WoW.

*Note*: Only tested in Windows 7 x64.

It supports the following websites:
  
  + [http://www.wowhead.com/](http://www.wowhead.com/) - retrieve from item sets, item comparisons, and transmog sets. Also, it can get morph codes from items, npcs, or spells.
  + [http://us.battle.net/wow/](http://us.battle.net/wow/) - retrieve from person's armory
  + [http://wowroleplaygear.com/](http://wowroleplaygear.com/) - transmog set
  + [http://mogmygear.com/](http://mogmygear.com/) - transmog set
  + [http://www.worldofwardrobes.net](http://www.worldofwardrobes.net) - transmog set
  + [http://www.wowmogging.com](http://www.wowmogging.com) - transmog set

It also will support any website that has wowhead/wowdb links to items in them. So just try it out on a transmog site and see if it works. If it doesn't work then let me know at [issues](https://github.com/kdar/morphgen/issues).

#### Usage

Just run morphgen.exe for the GUI.

#### CLI Usage

This won't work in windows unless you redirect the output using 1> file. E.g.: 

    morphgen <url> 1>output.txt

run `morphgen --help` for options

Wowhead:

    tmorphgen http://www.wowhead.com/itemset=523

    .item 1 22418
    .item 3 22419
    .item 5 22416
    .item 6 22422
    .item 7 22417
    .item 8 22420
    .item 9 22423
    .item 10 22421

Wow armory:

    tmorphgen http://us.battle.net/wow/en/character/tichondrius/Dominozx/advanced

    .item 1 32235
    .item 3 50646
    .item 4 41254
    .item 5 50649
    .item 6 50707
    .item 7 50696
    .item 8 50607
    .item 9 84972
    .item 10 50615
    .item 15 84804
    .item 16 39344
    .item 17 39344
    .enchant 1 4444
    .enchant 2 5035
    .face 7
    .gender 0
    .hair 9
    .haircolor 2
    .race 5
    .skin 5

Wow roleplay gear:

    tmorphgen http://wowroleplaygear.com/2010/01/11/strength-of-the-clefthoof/

    .item 3 15169
    .item 5 25689
    .item 6 9827
    .item 7 25690
    .item 8 25691
    .item 9 24700
    .item 10 15310
    .item 16 50319

Mog my gear:

    tmorphgen "http://mogmygear.com/gallery.php?r=2&g=f&c=1&q=2&p=1&s=9"

    .item 1 8317
    .item 3 8319
    .item 5 8312
    .item 6 8315
    .item 7 8318
    .item 8 8316
    .item 9 8311
    .item 10 8314

World of Wardrobes:

    tmorphgen http://www.worldofwardrobes.net/2011/10/bristlebark-leather-set/

    .item 3 14573
    .item 5 14570
    .item 6 14567
    .item 7 14574
    .item 8 14568
    .item 9 14569
    .item 10 14572
    .item 16 14571
    .item 17 15894



#### Compiling

I only tested this on Windows. I'm assuming little work would need to be done to get it to compile for OSX.

Used gcc (tdm64-1) 4.6.1 for windows, and downloaded Lua 5.1.4 source and built it from [http://luabinaries.sourceforge.net/download.html](http://luabinaries.sourceforge.net/download.html).

    build.bat

#### Lua

I included in the distribution [github.com/aarzilli/golua](github.com/aarzilli/golua), modified lua.go and added `#cgo windows pkg-config: lua5.1`, and added pkg-config.bat so it would complile on windows.

pkg-config.bat also has to be in the main respository, and to luar to compile.

