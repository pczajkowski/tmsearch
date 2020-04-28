let app = {};

app.filterTable = function(value) {
	let data = document.getElementById("dataRows");
	let rows = data.getElementsByTagName("tr");
	if (value === "") {
		for (let i = 0; i < rows.length; i++) {
			let row = rows[i];
			row.style.display = "";
		}
		return;
	}

	const filter = value.toUpperCase();	

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

app.sorted = " (sorted)";

app.clearFilter = function() {
	let filter = document.getElementById("filter");
	filter.value = "";
	app.filterTable("");

	document.querySelectorAll('th').forEach(th => {
		th.textContent = th.textContent.replace(app.sorted, "");
	});
};

app.getCellValue = (tr, idx) => tr.children[idx].innerText || tr.children[idx].textContent;

app.comparer = (idx, asc) => (a, b) => ((v1, v2) => 
		v1 !== '' && v2 !== '' && !isNaN(v1) && !isNaN(v2) ? v1 - v2 : v1.toString().localeCompare(v2)
	)(app.getCellValue(asc ? a : b, idx), app.getCellValue(asc ? b : a, idx));

app.sortTable = function() {
	// borrowed from https://stackoverflow.com/a/49041392/7342859
	document.querySelectorAll('th').forEach(th => th.addEventListener('click', (() => {
		const table = document.getElementById("dataRows");
		Array.from(table.querySelectorAll('tr'))
			.sort(app.comparer(Array.from(th.parentNode.children).indexOf(th), this.asc = !this.asc))
			.forEach(tr => table.appendChild(tr) );

		if (!th.textContent.includes(app.sorted)) {
			th.textContent += app.sorted;
		}

		th.parentNode.childNodes.forEach((child) => {
			if (child != th) {
				child.textContent = child.textContent.replace(app.sorted, "");
			}
		});
	})));

};
