require "iuplua"

module("about", package.seeall)

-- local L = package.loaded
-- print(L["ui"].VERSION)

function show(version) 
  iup.Message("About MorphGEN", [[Made by Kevin and Cassandra Darlington
Binary: http://github.com/kdar/morphgen-binary
Code: http://github.com/kdar/morphgen
Site: http://outroot.com

Version: ]] .. version)
end
