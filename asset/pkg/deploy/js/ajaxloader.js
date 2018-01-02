(function( $ ) {
    $.fn.loading = function () {

        // create loading element
        var loadingElement = document.createElement('div');
        loadingElement.id = 'ajaxloader';

        // apply styles
        loadingElement.style.position = 'absolute';
        loadingElement.style.margin = 'auto';
        loadingElement.style.top = 0;
        loadingElement.style.right = 0;
        loadingElement.style.bottom = 0;
        loadingElement.style.left = 0;
        loadingElement.style.background = 'url(/asset/img/ajax-loader.gif) no-repeat center center';
        loadingElement.style.width = '32px';
        loadingElement.style.height = '32px';
        loadingElement.style.display = 'none';

        // attach it to DOM
        $(this).append(loadingElement);

        // every time ajax is called
        $(document).ajaxSend(function () {
            $(loadingElement).show();
        })

        // every time ajax is completed
        $(document).ajaxStop(function () {
           $(loadingElement).hide();
        });
    };

})(jQuery);

$(document).ready(function () {

    $('body').loading();
});