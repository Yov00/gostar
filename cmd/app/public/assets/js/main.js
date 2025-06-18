
async function deleteUser(e) {
  const email = e.currentTarget.getAttribute("data-email")

  if (email != null && email.trim().length > 0) {
    const URL = '/delete/' + email;
    const res = await fetch(URL, {
      method: "Delete"
    })

    if (res.status == 200) {
      container = e.target.closest(".stats.shadow");
      if (container) {
        container.classList.add("removing")
        setTimeout(() => {
          container.remove()
        }, 280)
      }
    }

  }
}

async function addUser() {
  try {
    const res = await fetch("/addUser", {
      method: "POST",
      body: JSON.stringify({ name: "example", email: "example@test.com" }),
    })
    const data = await res.text()
    const userUl = document.querySelector("#user-list")
    if (userUl) {
      userUl.insertAdjacentHTML("beforeend", data)
    }
  } catch (err) {
    console.error(err)
  }
}
