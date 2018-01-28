// Switchery

// enable or disable application
function toggleApp(uuid, enabled){
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
                location.reload();
            }else{
                alertify.alert("Error toggle app!");
            }
        },
        error: function(xhr, status){
            alertify.alert("Error toggle app");
        },
    });
}

$(document).ready(function() {
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
            toggleApp(uuid, enabled);
        });
    }
});
// /Switchery