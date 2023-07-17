const formInput = document.querySelectorAll('.form-input')
const formField = document.querySelectorAll('.form-field')
formInput.forEach((item, i )=> {
  if (item.value != '') {
    formField[i].classList.add('active')
  }
  item.addEventListener("input", () => {
    formField[i].classList.remove('active')
    if (item.value != '') {
      formField[i].classList.add('active')
    }
  })
})