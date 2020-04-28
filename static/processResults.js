let app = {};

app.filterTable = function(value) {
	const filter = value.toUpperCase();
	let data = document.getElementById("dataRows");
	let rows = data.getElementsByTagName("tr");

	for (let i = 0; i < rows.length; i++) {
		let row = rows[i];
		row.style.display = "none";

		let cells = row.getElementsByTagName("td");
		for (let j = 0; j < cells.length; j++) {
			let cell = cells[j];

			if (cell.innerText.toUpperCase().indexOf(filter) > -1) {
				row.style.display = "";
				break;
			}
		}
	}
};

app.clearFilter = function() {
	let table = document.getElementsByTagName("tbody")[0];
	let tr = table.getElementsByTagName("tr");

	for (let i = 0; i < tr.length; i++) {
		tr[i].style.display = "";
	}

	let filter = document.getElementById("filter");
	filter.value = "";
};

app.sortTable = function() {
	// borrowed from https://stackoverflow.com/a/49041392/7342859
	const getCellValue = (tr, idx) => tr.children[idx].innerText || tr.children[idx].textContent;

	const comparer = (idx, asc) => (a, b) => ((v1, v2) => 
		v1 !== '' && v2 !== '' && !isNaN(v1) && !isNaN(v2) ? v1 - v2 : v1.toString().localeCompare(v2)
	)(getCellValue(asc ? a : b, idx), getCellValue(asc ? b : a, idx));

	// do the work...
	document.querySelectorAll('th').forEach(th => th.addEventListener('click', (() => {
		const table = th.closest('table');
		Array.from(table.querySelectorAll('tr:nth-child(n+2)'))
			.sort(comparer(Array.from(th.parentNode.children).indexOf(th), this.asc = !this.asc))
			.forEach(tr => table.appendChild(tr) );

		const sorted = " (sorted)";
		if (!th.textContent.includes(sorted)) {
			th.textContent += sorted;
		}

		th.parentNode.childNodes.forEach((child) => {
			if (child != th) {
				child.textContent = child.textContent.replace(sorted, "");
			}
		});
	})));

};
