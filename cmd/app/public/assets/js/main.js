
async function deleteUser(e) {
    const email = e.currentTarget.getAttribute("data-email")

    if (email != null && email.trim().length > 0) {
        const URL = '/delete/' + email;
        const res = await fetch(URL, {
            method: "Delete"
        })

       if(res.status == 200){
            container = e.target.closest(".stats.shadow");
            console.log(container)
            if(container){
                container.remove()
            }
       }

    }
}