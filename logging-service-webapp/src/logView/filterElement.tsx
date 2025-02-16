export enum FilterType {
    TEXT = "text",
    NUMBER = "number",
    DATETIME = "datetime-local"
}

// Map each FilterType to its expected value type
type FilterValueMap = {
    [FilterType.TEXT]: string;
    [FilterType.NUMBER]: number;
    [FilterType.DATETIME]: string; // datetime-local returns a string
};

interface FilterElementProps<T extends FilterType> {
    label: string;
    filterType: T;
    value: FilterValueMap[T];
    onChange: (value: any) => void;
}

const FilterElement = <T extends FilterType>({ label, filterType, value, onChange }: FilterElementProps<T>) => {
    return (
        <label className="border rounded-md p-1">
            {label}:
            <input
                className="border p-1 ml-1"
                type={filterType}
                value={value}
                onChange={(e) => onChange(e.target.value)}
            />
        </label>
    )
}

export default FilterElement;
