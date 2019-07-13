// Adds the results with buttons that will redirect to the obtain page.
var clickedbutton = false
function updateprofit() {


    var itemrecipeid = document.getElementsByClassName("itemrecipeid")
    for (var i = 0; i < itemrecipeid.length; i++) {
        var recipeid = itemrecipeid.item(i).innerHTML

        const updatebutton = (
            <button class="uk-button uk-button-default updatebutton" value={recipeid}>Update</button>
        )


        ReactDOM.render(updatebutton, document.getElementById(recipeid))
        // We're going to have to make the updatebutton onclick work.
        document.getElementsByClassName("updatebutton").item(i).onclick = function () {
            if (!clickedbutton) {
                clickedbutton = true
                obtainrecipe(this.value)
                // We don't ever unlock the button.
                // So that it doesn't interrupt the current inserts and api calls.
            }
        }
    }

}


function obtainrecipe(recipeid) {
    $("#obtaininput").attr("value", recipeid)
    document.getElementById("obtainform").submit();
}


function unlock() {
    clickedbutton = false
}