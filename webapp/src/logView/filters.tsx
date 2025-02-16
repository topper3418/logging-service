
import React from "react";
import { LogFilters } from "./logHooks";
import FilterElement, { FilterType } from "./filterElement";

interface LogFiltersProps {
    logFilters: LogFilters;
}

const Filters: React.FC<LogFiltersProps> = ({ logFilters: { get, set, clear } }) => {

    return (
        <div className="flex gap-2.5">
            <FilterElement
                filterType={FilterType.DATETIME}
                label="Min Time"
                value={get.minTime}
                onChange={set.minTime} />
            <FilterElement
                filterType={FilterType.DATETIME}
                label="Max Time"
                value={get.maxTime}
                onChange={set.maxTime} />
            <FilterElement
                filterType={FilterType.NUMBER}
                label="Offset"
                value={get.offset}
                onChange={set.offset} />
            <FilterElement
                filterType={FilterType.NUMBER}
                label="Limit"
                value={get.limit}
                onChange={set.limit} />
            <FilterElement
                filterType={FilterType.TEXT}
                label="Search"
                value={get.search}
                onChange={set.search} />
            <button onClick={clear}>
                Clear
            </button>
        </div>
    )
}

export default Filters;
