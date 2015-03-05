require "iuplua"

module("about", package.seeall)

-- local L = package.loaded
-- print(L["ui"].VERSION)

function show(version) 
  iup.Message("About MorphGEN", [[Made by Kevin and Cassandra Darlington
Site: http://outroot.com

Version: ]] .. version)
end
