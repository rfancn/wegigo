function validateSteps() {
    var errorSteps = [];
    var totalSteps = 4;
    for(var i=0;i<totalSteps;i++){
        var ret = runFunction("validateStep"+i);
        if(typeof(ret) !="undefined" && ret.length != 0) {
            errorSteps.push(i);
        }
    }
    if(errorSteps.length > 0){
        //$('#smartwizard').smartWizard("stepState", errorSteps, "error");
        alertify.alert("You need correct the errors in step" + errorSteps+":<br/>"+ret);
        return false;
    }

    return true
}

function onClickFinish(){
    if(!validateSteps()) {
        return false
    }

    $.ajax({
        url: "/deploy/config",
        dataType: 'json',
        contentType:"application/json; charset=utf-8",
        type: 'POST',
        data: JSON.stringify({
            'general': vueGeneral.$data,
            'wechat': vueWechat.$data,
            'server': tableServers.rows().data().toArray(),
        }),
        success: function(data) {
            switch(data.Result) {
                case "error":
                    alertify.alert(data.Detail);
                    break;
                case "success":
                    window.location.replace(data.Detail);
                    break;
            };
        },
        error: function(xhr, status){
            console.log("config error");
            alertify.alert("Http Error:"+status);
        },

    });

    return true;
}

function runFunction(name, arguments)
{
    var fn = window[name];
    if(typeof fn !== 'function')
        return;

    return fn.apply(window, arguments);
}

function init_smartwizard(){
    $('#smartwizard').smartWizard({
        keyNavigation:false, // Enable/Disable keyboard navigation(left and right keys are used if enabled)
        toolbarSettings: {
            toolbarExtraButtons: [
                $('<button id="btn-finish"></button>').text('Finish').addClass('btn btn-primary').on('click', function(){
                    onClickFinish();
                }),
            ],
        },
    });

    $("#btn-finish").hide();

    // Initialize the leaveStep event
    $("#smartwizard").on("leaveStep", function(e, anchorObject, stepNumber, stepDirection) {
        return runFunction("leaveStep"+stepNumber, [stepDirection]);
    });

    $("#smartwizard").on("showStep", function(e, anchorObject, stepNumber, stepDirection) {
        if(stepNumber == 3) {
            $("#btn-finish").show();
        }else{
            $("#btn-finish").hide();
        }
        return runFunction("showStep"+stepNumber, [stepDirection]);
    });
}

$(document).ready(function(){
    init_smartwizard();
    //setup ajax loader
    $(document).ajaxStart(function(){
        $.LoadingOverlay("show");
    });
    $(document).ajaxStop(function(){
        $.LoadingOverlay("hide");
    });
});
