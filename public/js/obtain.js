
function removeemptymaterials() {

    $("#materiallist li[value='0']").remove();

    $("#mainrecipelist li[value='0']").remove();

    $("#materialinfo li[value='0']").remove();
}

// Uses .itemicon values, and converts the src to the actual image location.
function getitemicon() {
    var x = $(".itemicon");
    for (var i = 0; i < x.length; i++) {
        var srcvalue = '0' + x.eq(i).attr("value");
        var folderidentifier = srcvalue.substring(1, 3);
        var folder = '0' + folderidentifier + '000/';
        var path = '/public/img/icon/' + folder + srcvalue + '.png';
        x.eq(i).attr("src", path);
    }

}

// Grabs from Obtain, the .itemamount, and .totalcost, to give us the real total costs of individual items.
function gettotalcosts() {
    var itemamount = $(".itemamount");
    var totalcost = $(".totalcost");
    for (var i = 0; i < itemamount.length; i++) {
        var tempstring = 'Total Cost: ' + itemamount.eq(i).attr("value") * totalcost.eq(i).attr("value");
        totalcost.eq(i).html(tempstring);
    }
}

// Changes all .pricenumbers to comma'd numbers
function changetodecimals() {
    var pricenumber = $(".pricenumber")
    for (var i = 0; i < pricenumber.length; i++) {
        var newnumber = pricenumber.eq(i).html().toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
        pricenumber.eq(i).html(newnumber)
    }

}

function changeunixtodate() {
    var unixtime = $(".datetime")
    for (var i = 0; i < unixtime.length; i++) {
        var newtime = parseInt(unixtime.eq(i).attr("value")) * 1000
        var date = new Date(newtime)
        var days = date.getDate()
        var months = date.getMonth()
        var years = date.getFullYear()
        var hours = date.getHours()
        var minutes = date.getMinutes()
        if (hours >= 12) {
            var ampm = "pm"
        } else {
            var ampm = "am"
        }
        hours = hours % 12
        if (hours == 0) {
            hours = 12
        }
        if (minutes < 10) {
            minutes = "0" + minutes
        }
        var datestring = "Added: " + days + "/" + months + "/" + years + " at " + hours + ":" + minutes + " " + ampm

        unixtime.eq(i).html(datestring)
    }
}

function getavatar() {
    var user = $("#id").attr("value")
    var avatarhash = $("#avatar").attr("value")
    var path = "https://cdn.discordapp.com/avatars/" + user + "/" + avatarhash + ".gif"

    var userimage = $("#userimage")
    userimage.attr("src", path)
    $("#discorduser").html($("#username").attr("value") + "#" + $("#discriminator").attr("value"))
}

// Instead of iterating through golang's templates everytime, we just the map results,
// and create the document lists through javascript itself.
function innermaterials() {

    var innermatmap = $(".innermatmaps")
    for (var i = 0; i < innermatmap.length; i++) {
        if (innermatmap.eq(i).attr("value") == "[]") {
            innermatmap.eq(i).remove()
        } else {
            var divaccordian = document.createElement("div")
            divaccordian.className = "uk-accordion-content"

        }
    }
}