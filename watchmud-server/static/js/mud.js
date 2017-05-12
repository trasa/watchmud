;(function($) {


    $.mudclient = function (element, options) {

        var defaults = {
            foo: 'bar',
            socketUrl: "ws://localhost:8888/ws",

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
        }

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
                // var msg = JSON.parse(evt.data);
                // handleMessage(msg);
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

                function displayMessage(msg) {
                    var display = $('#display');
                    display
                        .append(msg + "\n\n")
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