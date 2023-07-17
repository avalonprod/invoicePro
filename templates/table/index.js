const tableBody = document.querySelector(".table-body") 


for (let i = 50; i >= 0; i--) {
  tableBody.innerHTML += `
    <tr>
      <td class="td-invoice-name">Invoice</td>
      <td class="td-client-name">Invoice</td>
      <td class="td-lable-name">Invoice</td>
      <td class="td-date">Invoice</td>
      <td class="td-balance">Invoice${i}</td>
     </tr>
  `
}