_G.callbacks = _G.callbacks or {}

require "iuplua"
require "ui/about"

module("ui", package.seeall)

local cancelProgress
local updateBar

-- Actions

function do_about()
  about.show(VERSION)
  return iup.DEFAULT
end

function do_close()
  return iup.CLOSE
end

function do_generate()
  statusLabel.title = "Generating TMorph codes. Could take a while..."
  btn1.active = "NO"
  output.value = ""
  progressBar.visible = "YES"

  generate({
    url=textUrl.value,
    notmog=optionTmog.value=="OFF",
  })
 
  return iup.DEFAULT
end

function do_download()
  err = download()
  if err ~= nil then
    show_error(err)
  end
end

-- Callbacks

-- called when generate is done
function _G.generate_callback(data, err)
  if err ~= nil then    
    show_error(err)
  else
    output.value = data
    statusLabel.title = "Done."
  end  
  btn1.active = "YES"
  progressBar.visible = "NO"  
  cancelProgress = true
end

function _G.checkupdate_callback(update, err)
  if err ~= nil then
    statusLabel.title = "Updates: " .. err
  else
    if update ~= "" then
      statusLabel.title = "Ready."
      show_updatebar(update)
    else
      statusLabel.title = "You have the latest version! Ready."
    end
  end
end

-- UI functions

function show_updatebar(text)
  updateLabel = iup.label{
    expand="horizontal",
    title=text,
    ellipsis="YES",
    visible="YES",
  }

  updateButton = iup.button{
    title="Download", 
    size="80x", 
    action=do_download,
    visible="YES",
  }

  updateBar = iup.hbox{
    margin="3x3",
    alignment="ACENTER";
    updateLabel, iup.fill{},
    updateButton,
    visible="YES",
  }

  iup.Insert(vbox, statusBar, updateBar)
  updateBar:map()
  dlg:map()
end

function hide_updatebar()
  updateBar:destroy()
end

function show_error(err)
  output.value = err
  tags = iup.user { bulk = "Yes", cleanout = "Yes" }
  iup.Append(tags, iup.user { selectionpos = "0:1000", fgcolor = "255 0 0"})
  output.addformattag = tags    
  statusLabel.title = "Error occurred."
end

-- Menu

mmenu = {
  "File",{
    "Exit",do_close,
  },
  "Help",{
    "About",do_about,
  }
}

function create_menu(templ)
  local items = {}
  for i = 1,#templ,2 do
    local label = templ[i]
    local data = templ[i+1]
    if type(data) == 'function' then
      item = iup.item{title = label}
      item.action = data
    elseif type(data) == 'nil' then
      item = iup.separator{}
    else
      item = iup.submenu{create_menu(data); title = label}
    end
    table.insert(items,item)
  end
  return iup.menu(items)
end

menu = create_menu(mmenu)

-- Body

btn1 = iup.button{
  title="Generate", 
  size="80x", 
  action=do_generate
}

textUrl = iup.text{
  value="", 
  expand="HORIZONTAL", 
  tip="Put in a URL to a website to generate codes from. Could be wowhead, wowarmory, or a number of transmog sites."
}

bbox = iup.hbox{
  textUrl,btn1; 
  gap=4; 
  margin="4x4"
} 
output = iup.text{
  multiline="YES",
  expand="YES",
  formatting="YES",
  readonly="YES",
  wordwrap="YES",
  tip="The output of the generation",
} 

optionTmog = iup.toggle{
  title="Use transmogged items (armory only)", 
  value="ON", 
  tip="If checked, it will use the character's transmogrified items on the armory.",
}
optionsBox = iup.hbox{
  optionTmog; 
  gap=4; 
  margin="4x4"
}

progressBar = iup.progressbar{
  rastersize="60x20",
  marquee="YES",
  visible="NO",
  value=0.0,
}

statusLabel = iup.label{
  expand="horizontal",
  title="Ready.",
  ellipsis="YES",
}

statusBar = iup.hbox{
  margin="3x3",
  alignment="ACENTER";
  statusLabel, iup.fill{},
  progressBar
}

vbox = iup.vbox{
  bbox,
  optionsBox,
  output,
  iup.label{separator="horizontal"},
  statusBar,
}

dlg = iup.dialog{
  vbox; 
  title="MorphGEN", 
  size="350x230",
  minsize="350x320", 
  shrink="YES",
  menu=menu, 
  icon="resources/icon.ico"
}
dlg.defaultenter = btn1

-- our callback processor.
-- we do this because GUI draws can only
-- be done in the main thread. so when we call
-- a function in Go that needs to do a lot of processing,
-- it can do a goroutine and then add a callback later
-- once it's done.
timer = iup.timer{time=100}
function timer:action_cb()
  for k, v in pairs(_G.callbacks) do
    _G.callbacks[k] = nil 
    _G[k](unpack(luar.slice2table(v)))
  end
  _G.callbacks = {}
end
timer.run = "YES"

-- check for updates
timeronce = iup.timer{time=500}
function timeronce:action_cb()
  timeronce.run = "NO"
  statusLabel.title = "Checking for updates..."
  checkupdate()
end
timeronce.run = "YES"

dlg:show()
iup.MainLoop()
