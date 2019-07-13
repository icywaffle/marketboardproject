// Adds the results with buttons that will redirect to the obtain page.
var clickedbutton = false
function updateprofit() {


    var itemrecipeid = document.getElementsByClassName("itemrecipeid")
    var checktime = document.getElementsByClassName("datetime")
    for (var i = 0; i < itemrecipeid.length; i++) {
        var now = new Date
        now = now.getTime() / 1000
        var profitupdatetime = checktime.item(i).getAttribute("value")
        var timesinceupdate = now - profitupdatetime
        // Using 87000 instead 86400 as a better buffer to update.
        // We only show the button if it's updateable.
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
        if (timesinceupdate > (87000 / 2)) {
        } else {
            // Just hide the button for the user, since there's a cooldown anyway
            itemrecipeid.item(i).style.display = "none"
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