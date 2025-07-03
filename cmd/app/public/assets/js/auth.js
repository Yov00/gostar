console.log('auth loaded')

function register_submit(e) {
  e.preventDefault()
  const name = e.target.querySelector("#name")?.value;
  const email = e.target.querySelector("#email")?.value;
  const password = e.target.querySelector("#password")?.value;
  const confirmPassword = e.target.querySelector("#confirmPassword").value;

  if (email?.trim().length === 0) {
    alert("Email is required")
    return
  }

  if (name?.trim().length === 0) {
    alert("Name is required")
    return
  }

  if (password?.trim().length < 6) {
    alert("Password length should be at least 6 characters")
    return
  }

  if (password !== confirmPassword) {
    alert("Passwords don' match")
    return
  }


  alert("send success!")
  e.target.submit()
}
