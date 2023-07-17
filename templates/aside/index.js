const aside = document.querySelector(".aside")
const asideToggler = document.querySelector(".aside-toggler")

asideToggler.addEventListener("click", () => {
  aside.classList.toggle("active")
})