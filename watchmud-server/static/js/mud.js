;(function($) {


    $.mudclient = function(element, options) {
        console.log("hi");


        // to avoid confusions, use "plugin" to reference the
        // current instance of the object
        var plugin = this;


        // ctor
        plugin.init = function() {
            console.log("init");
        };




        $.fn.mudclient = function (options) {
            return this.each(function () {
                if (undefined == $(this).data('mudclient')) {
                    var plugin = new $.mudclient(this, options);
                    $(this).data('mudclient', plugin);
                }
            });
        };

        plugin.init();
    };
    
    
    var app = $.sammy(function() {
        this.use(Sammy.EJS);

        this.get('#/', function() {
            this.render('templates/index.ejs', function(html) {
                $('#mainContent').html(html);
                
                var mc = $.mudclient( ); /*.data('mudclient');*/
                displayMessage("Connecting to socket...");
                // mc.run();
                
                
                
                function displayMessage(msg) {
                    var display = $('#display');
                    display
                        .append(msg + "\n\n")
                        .stop()
                        .animate({ scrollTop: display[0].scrollHeight}, 800);
                }
            });
        });
    });

    $(function() {
        app.run('#/');
    });
})(jQuery);