
import React from "react";
import { LogFilters } from "./logHooks";
import FilterElement, { FilterType } from "./filterElement";

interface LogFiltersProps {
    logFilters: LogFilters;
}

const Pagination: React.FC<LogFiltersProps> = ({ logFilters: { get, set } }) => {
    const pageNumber = Math.floor(get.offset / get.limit) + 1;
    const nextPage = () => set.offset(get.offset + get.limit);
    const prevPage = () => set.offset(get.offset - get.limit);

    return (
        <div className="flex gap-2.5">
            <FilterElement
                filterType={FilterType.NUMBER}
                label="Limit"
                value={get.limit}
                onChange={set.limit} />
            <div className="flex gap-1 border rounded-md">
                <button onClick={prevPage}>&#x25C0;</button>
                <div className="flex flex-col justify-center m-0 p-0 w-15">
                    <p className="text-center">
                        {pageNumber}
                    </p>
                </div>
                <button onClick={nextPage}>&#x25B6;</button>
            </div>
        </div>
    )
}

export default Pagination;
