
function removeemptymaterials() {

    $("#materiallist li[value='0']").remove();

}


function getitemicon() {
    var srcvalue = '0' + $(".itemicon").attr("src");
    var folderidentifier = srcvalue.substring(1, 3);
    var folder = '0' + folderidentifier + '000/';
    var path = '/public/img/icon/' + folder + srcvalue + '.png';
    $(".itemicon").attr("src", path);
}