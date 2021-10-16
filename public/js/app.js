function actionCellRenderer(params) {
    let eGui = document.createElement("div");

    let editingCells = params.api.getEditingCells();
    // checks if the rowIndex matches in at least one of the editing cells
    let isCurrentRowEditing = editingCells.some((cell) => {
        return cell.rowIndex === params.node.rowIndex;
    });

    if (isCurrentRowEditing) {
        eGui.innerHTML = `
            <button
            class="action-button update"
            data-action="update">
                update
            </button>
            <button
            class="action-button cancel"
            data-action="cancel">
                cancel
            </button>
            `;
    } else {
        eGui.innerHTML = `
            <button
            class="action-button edit"
            data-action="edit">
                edit
            </button>
            <button
            class="action-button delete"
            data-action="delete">
                delete
            </button>
            `;
    }

    return eGui;
}


$(document).ready(function() {
    const columnDefs = [
        { colId: "ID", field: "ID", sortable: true, hide: true, resizable: true },
        { field: "country", width: 220, sortable: true, filter: true, resizable: true, editable: true },
        { field: "operator", width: 330, sortable: true, filter: true, resizable: true, editable: true },
        { field: "code.mcc", width: 75, headerName: "MCC", sortable: true, filter: true, resizable: true, editable: true },
        { field: "code.mnc", width: 75, headerName: "MNC", sortable: true, filter: true, resizable: true, editable: true },
        { headerName: "action", minWidth: 150, cellRenderer: actionCellRenderer, editable: false, colId: "action" },
    ];


    const gridOptions = {
        columnDefs: columnDefs,

        enableCellTextSelection: true,
        onModelUpdated : function () {
            this.api.sizeColumnsToFit();
        },
        suppressClickEdit: true,
        onCellClicked(params) {
            // Handle click event for action cells
            if (params.column.colId === "action" && params.event.target.dataset.action) {
                let action = params.event.target.dataset.action;

                if (action === "edit") {
                    params.api.startEditingCell({
                    rowIndex: params.node.rowIndex,
                    // gets the first columnKey
                    colKey: params.columnApi.getDisplayedCenterColumns()[0].colId
                    });
                }

                if (action === "delete") {
                    params.api.applyTransaction({
                    remove: [params.node.data]
                    });
                }

                if (action === "update") {
                    params.api.stopEditing(false);
                }

                if (action === "cancel") {
                    params.api.stopEditing(true);
                }
            }
     },

    onRowEditingStarted: (params) => {
        params.api.refreshCells({
        columns: ["action"],
        rowNodes: [params.node],
        force: true
        });
    },
    onRowEditingStopped: (params) => {
        params.api.refreshCells({
        columns: ["action"],
        rowNodes: [params.node],
        force: true
        });
    },
    onCellValueChanged: function(event) {
        console.log('cellValueChanged', event);

        $.ajax({
            url: "/e212api.v1/e212/update",
            type: "POST",
            contentType: "application/json; charset=utf-8",
            traditional: true,
            data: JSON.stringify(event.data),

            success: function(data, status) {
                console.log(`${data} and status ${status}`)
            },
            error: function(errMsg) {
                alert(errMsg);
                event.node.data[event.colDef.field] = event.oldValue;
                event.api.refreshCells({ rowNodes: [event.node], columns: [event.column.colId] });
            }
        })
    },
    editType: "fullRow",
};


    let  myGrid = $('#myGrid').get(0)

    agGrid.simpleHttpRequest({url: 'e212api.v1/e212/'})
    .then(data => {
        gridOptions.api.setRowData(data);
    });

    let grid = new agGrid.Grid(myGrid, gridOptions);

    $(".search").keyup(function () {
        var searchTerm = $(".search").val();
        gridOptions.api.setQuickFilter(searchTerm);

    });
  });
