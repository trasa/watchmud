;(function($) {


    $.mudclient = function (element, options) {

        var defaults = {
            foo: 'bar',
            socketUrl: "ws://" + window.location.host + "/ws",

            // if your plugin is event-driven, you may provide callback capabilities
            // for its events. execute these functions before or after events of your
            // plugin, so that users may customize those particular events without
            // changing the plugin's code
            onFoo: function () {
            }


        };

        // to avoid confusions, use "plugin" to reference the
        // current instance of the object
        var plugin = this;

        // this will hold the merged default, and user-provided options
        // plugin's properties will be available through this object like:
        // plugin.settings.propertyName from inside the plugin or
        // element.data('pluginName').settings.propertyName from outside the plugin,
        // where "element" is the element the plugin is attached to;
        plugin.settings = {};


        var $element = $(element), // reference to the jQuery version of DOM element
            element = element;    // reference to the actual DOM element

        var websocket;

        // ctor
        plugin.init = function () {
            plugin.settings = $.extend({}, defaults, options);
            console.log("init");
        };

        /* private */
        var displayMessage = function (msg) {
            plugin.settings["displayMessage"](msg);
        };

        var displayError = function(msg) {
            displayMessage("Error: " + JSON.stringify(msg))
        };

        /* private */
        var handleMessage = function(msg) {
            switch(msg["msg_type"]) {

                case "enter_room":
                    handleEnterRoomNotification(msg);
                break;

                case "login_response":
                    handleLoginResponse(msg);
                    break;

                case "leave_room":
                    handleLeaveRoomNotification(msg);
                    break;

                case "look":
                    handleLookResponse(msg);
                    break;

                case "move":
                    handleMoveResponse(msg);
                    break;

                case "tell":
                    handleTellResponse(msg);
                    break;

                case "tell_notification":
                    handleTellNotification(msg);
                    break;

                case "tell_all":
                    handleTellAllResponse(msg);
                    break;

                case "tell_all_notification":
                    handleTellAllNotification(msg);
                    break;

                default:
                    displayMessage("Unknown message received: " + JSON.stringify(msg));
                    break;
            }
        };

        var handleEnterRoomNotification = function(msg) {
            displayMessage(msg["player"] + " enters.");
        };

        var handleLeaveRoomNotification = function(msg) {
            displayMessage(msg["player"] + " left.");
        };

        var handleLoginResponse = function(msg) {
            displayMessage("Login Response: Success=" + msg["success"] + " " + msg["result_code"]);
            displayMessage("Player is: " + JSON.stringify(msg["player"]));
        };

        var handleLookResponse = function(msg) {
            if (msg["success"]) {
                displayRoom(msg);
            } else {
                displayError(msg);
            }
        };

        var handleMoveResponse = function(msg) {
            if (msg["success"]) {
                displayRoom(msg);
            } else if (msg["result_code"] === "CANT_GO_THAT_WAY") {
                displayMessage("There's no exit that way.");
            } else {
                displayError(msg);
            }
        };

        var handleTellResponse = function(msg) {
            if (msg["success"]) {
                displayMessage("sent.");
            } else {
                displayMessage("tell failed: " + msg["result_code"])
            }
        };

        var handleTellNotification = function(msg) {
            displayMessage(msg["sender"] + " says '" + msg["value"] + "'");
        };

        var handleTellAllResponse = function(msg) {
            if (msg["success"]) {
                displayMessage("sent.");
            } else {
                displayMessage("Tell All failed: " + msg["result_code"]);
            }
        };

        var handleTellAllNotification = function(msg) {
            displayMessage("tell_all " + msg["sender"] +  "> " + msg["value"]);
        };


        var displayRoom = function(msg) {
            displayMessage(msg["room_name"]);
            displayMessage(msg["description"]);
            displayMessage("Exits: " + formatExits(msg["exits"]) + "\n")
        };

        var formatExits = function(exitStr) {
            if (exitStr === "") {
                return "None";
            }
            var s = "";
            for (var i = 0, len = exitStr.length; i < len; i++) {
                // alert(str[i]);
                switch(exitStr.charAt(i))  {
                    case "n":
                        s += "North, ";
                        break;
                    case "s":
                        s += "South, ";
                        break;
                    case "e":
                        s += "East, ";
                        break;
                    case "w":
                        s += "West, ";
                        break;
                    case "u":
                        s += "Up, ";
                        break;
                    case "d":
                        s += "Down, ";
                        break;
                }
            }
            if (s.length > 0) {
                s = s.substr(0, s.length -  2);
            }
            return s;
        };
        /* public */
        plugin.run = function () {
            displayMessage("Connecting to " + plugin.settings["socketUrl"]);
            websocket = new WebSocket(plugin.settings["socketUrl"]);

            websocket.onopen = function (evt) {
                displayMessage("opening socket...ready");
            };
            websocket.onclose = function (evt) {
                displayMessage("Socket closed");
                websocket = null;
            };
            websocket.onmessage = function (evt) {
                console.log("msg rec: ");
                console.log(evt.data);
                // TODO
                var msg = JSON.parse(evt.data);
                handleMessage(msg);
            };
            websocket.onerror = function (evt) {
                displayMessage("Error Received");
                console.log("Error:");
                console.log(evt);
            };
        };

        /* public */
        plugin.send = function(msg) {
            msg = JSON.stringify(msg);
            console.log("send " + msg);
            websocket.send(msg);
        };

        /* public */
        // plugin.forceClose = function() {
        //     console.log("force close!");
        //     websocket.close();
        // };

        plugin.init();
    };


    $.fn.mudclient = function (options) {
        return this.each(function () {
            if (undefined == $(this).data('mudclient')) {
                var plugin = new $.mudclient(this, options);
                $(this).data('mudclient', plugin);
            }
        });
    };

    var app = $.sammy(function () {
        this.use(Sammy.EJS);

        this.get('#/', function () {
            this.render('templates/index.ejs', function (html) {
                $('#mainContent').html(html);
                var mc;

                $('#connect').click(function () {
                    mc = $('#connect').mudclient({
                        displayMessage: displayMessage
                    }).data('mudclient');
                    mc.run();
                });

                $('#login').click(function() {
                   mc.send({
                       "msg_type" : "login",
                       "player_name" : $('#player').val(),
                       "password" : $('#pass').val()
                   });
                });

                $('#buf').keypress(function(event){
                    var keycode = (event.keyCode ? event.keyCode : event.which);
                    if(keycode === 13){
                        doSend();
                    }
                });

                function doSend() {
                    var buf = $('#buf');
                    mc.send({
                        "format":"line",
                        "value" : buf.val()
                    });
                    buf.val("");
                }

                $('#send').click(function() {
                   doSend();
                });


                function displayMessage(msg) {
                    var display = $('#display');
                    display
                        .append(msg + "\n")
                        .stop()
                        .animate({scrollTop: display[0].scrollHeight}, 800);
                }
            });
        });
    });

    $(function () {
        app.run('#/');
    });
})(jQuery);