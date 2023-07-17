const formInput = document.querySelectorAll('.form-input')
const formField = document.querySelectorAll('.invoice-item-field')
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


const showDetails = document.querySelector('.show-details')
const detailsWrapper = document.querySelector('.details-wrapper')

showDetails.addEventListener('click', ()=>{
  if(detailsWrapper.classList.contains('active')){
    detailsWrapper.classList.remove('active')
    showDetails.textContent = 'Show additional business details'
  } else{
    detailsWrapper.classList.add('active')
    showDetails.textContent = 'Hide additional business details'
  }
})