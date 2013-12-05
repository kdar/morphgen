var ws = null;
var stack_bar_bottom = {"dir1": "up", "dir2": "right", "spacing1": 0, "spacing2": 0};

var maskOptions = {
   spinner: {
     color: "white",
     lines: 10, length: 4, width: 2, radius: 5
   }
};

function notify(msg, type) {
  var opts = {
    title: msg,
    addclass: "stack-bar-bottom",
    cornerclass: "",
    width: "70%",
    type: type,
    sticker: false,
    mouse_reset: false,
    delay: 3000,
    stack: stack_bar_bottom
  };
  $.pnotify(opts);
}

function msg(name, args) {
  if (ws !== null) {
    ws.send(JSON.stringify({
      "Name": name,
      "Args": args
    }));
  }
};

window.checkupdate_callback = function(args) {
  if (args[1] != null) {
    notify(args[1], "error");
  } else {
    if (args[0] == "") {
      // notify("You have the latest version installed!", "info");
    } else {
      notify(args[0], "info");
    }
  }
};

window.generate_callback = function(args) {
  $('body').unmask();
  $("#generate").removeProp("disabled");
  // $(".progress > div").css("width", "0px");

  if (args[1] != null) {
    console.log(args);
    notify(args[1], "error");
  } else {
    $("#output").val(args[0]);
  }
};

// function onWindowResize() {
//   var bodyheight = $(document).height();
//   $("#container").height(bodyheight-12);
// };

$(function() {
  // onWindowResize();
  // $(window).resize(onWindowResize);

  $('body').mask(maskOptions);

  ws = new WebSocket("ws://" + location.host + "/ws");
  //ws.binaryType = "blob";
  ws.onopen = function() {
    $('body').unmask();
    msg("checkupdate");
  };
  ws.onmessage = function(e) {
    console.log(e.data);
    var data = JSON.parse(e.data);
    window[data.Name](data.Args);
  };
  ws.onclose = function(e) {
    ws = null;
    window.close();
  };

  $("body").bind("contextmenu", function(e) {
    if ($(e.target).is("textarea") || $(e.target).is("input")) {
      return true;
    }
    return false;
  });

  $("#generate").click(function() {
    $('body').mask(maskOptions);
    // $(".progress > div").css("width", "100%");
    $("#generate").prop("disabled", true);

    msg("generate", {"url": $("#url").val(), "notmog": $("#notmog").is(":checked")});
  });
});