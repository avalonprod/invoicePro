const tableBody = document.querySelector(".table-body") 


for (let i = 50; i >= 0; i--) {
  tableBody.innerHTML += `
<tr>
  <td>
    <input type="checkbox" class="marked-checkbox" id="marked-checkbox${i}">
    <label class="marked-label" for="marked-checkbox${i}"></label>
  </td>
  <td>Inv #127</td>
  <td>John Smith</td>
  <td>Jun 12.20.2024</td>
  <td>Wed</td>
  <td>$12343.32</td>
</tr>
  `
}