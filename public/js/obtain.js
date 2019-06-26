
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