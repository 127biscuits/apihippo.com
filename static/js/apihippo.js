// apihippo.js


    $('.verb').mouseover(function() {
        $('.hippo').css("color", "rgb(100, 100, 100)");
        $('.details').show();
    });
    $('.verb').mouseout(function() {
        $('.hippo').css("color", "rgb(199, 196, 196)");
        $('.details').hide();
    });
