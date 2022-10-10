function editPost(id) {
    var x = document.getElementById("post" + id)
    var y = document.getElementById("edit" + id)
    
    x.hidden = !x.hidden
    y.hidden = !y.hidden
}

function editOP() {
    var x = document.getElementById("threadpost")
    var y = document.getElementById("editthreadpost")

    x.hidden = !x.hidden
    y.hidden = !y.hidden
}