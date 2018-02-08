// enable or disable application
function toggleApp(uuid, enabled){
    // Switchery
    var toggleAppResult = false;

    $.ajax({
        url: "/wxmp/admin/app/toggle/" + uuid + "/" ,
        type: 'POST',
        data: {
            'enabled': enabled,
        },
        success: function(data) {
            console.log(data);
            if (data == "success") {
                console.log("in success");
                toggleAppResult = true;
                location.reload();
            }else{
                toggleAppResult = false;
                alertify.alert("Error toggle app!");
            }
        },
        error: function(xhr, status){
            toggleAppResult = false;
            alertify.alert("Error toggle app");
        },
    });

    return toggleAppResult
}

$(document).ready(function() {
    //setup ajax loader
    $(document).ajaxStart(function(){
        $.LoadingOverlay("show");
    });
    $(document).ajaxStop(function(){
        $.LoadingOverlay("hide");
    });

    if ($(".js-switch")[0]) {
        var elems = Array.prototype.slice.call(document.querySelectorAll('.js-switch'));
        elems.forEach(function (html) {
            var switchery = new Switchery(html, {
                color: '#26B99A'
            });
        });

        $(".js-switch").on("click", function(){
            var uuid =  $(this).prop('name');
            var enabled = $(this).prop('checked');
            return toggleApp(uuid, enabled);
        });
    }


});
// /Switchery