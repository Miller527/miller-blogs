$('#addRow').on( 'click', function () {
    $.ajax({
        url: "curd/role/add",
        type: "post",
        data: {
            "name":name
        },
        success: function (data) {
            console.log(data);
            //
            // if (data.user) {//验证成功
            //     location.href = "/index/"//location.href="/xxx/" jquery实现跳转
            // } else {//验证失败
            //     $('.error').text(data.msg).css({"color": "red", "margin-left": "10px"});//拿到登陆旁边的span进行提示
            //     //一秒后清空提醒
            //     setTimeout(function () {//定时器
            //         $(".error").text("")//由data.msg变成空字符串
            //     }, 1000)
            // }
        }
    })
    // t.row.add( [
    //     counter +'.1',
    //     counter +'.2',
    //     counter +'.3',
    //     v,
    // ] ).draw( false );
    //
    // counter++;
} );