
function removeemptymaterials() {

    $("#materiallist li[value='0']").remove();

    $("#mainrecipelist li[value='0']").remove();
}


function getitemicon() {
    var x = $(".itemicon");
    for (var i = 0; i < x.length; i++) {
        var srcvalue = '0' + x.eq(i).attr("src");
        var folderidentifier = srcvalue.substring(1, 3);
        var folder = '0' + folderidentifier + '000/';
        var path = '/public/img/icon/' + folder + srcvalue + '.png';
        x.eq(i).attr("src", path);
    }

}

function gettotalcosts() {
    var itemamount = $(".itemamount");
    var totalcost = $(".totalcost");
    for (var i = 0; i < itemamount.length; i++) {
        // There's no point in replacing the value
        if (itemamount.eq(i).attr("value") > 1) {
            var tempstring = 'Total Cost: ' + itemamount.eq(i).attr("value") * totalcost.eq(i).attr("value");
            totalcost.eq(i).html(tempstring);
        }
    }
}

function changetodecimals() {
    var pricenumber = $(".pricenumber")
    for (var i = 0; i < pricenumber.length; i++) {
        var newnumber = pricenumber.eq(i).html().toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
        pricenumber.eq(i).html(newnumber)
    }

}