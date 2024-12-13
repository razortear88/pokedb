function hellow(){
    $.ajax({url: "/pokemon/api/list", success: function(result){
        console.log(result)
        // $("#div1").html(result);
    }});
}