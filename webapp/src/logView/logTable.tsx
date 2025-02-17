import { CSSProperties } from "react";
import { LogFilters, LogsApi } from "./logHooks";
import Pagination from "./pagination";

export const formatDateString = (dateTimeString: string, includeDay: boolean = true) => {
    // Convert string to Date object
    const date = new Date(dateTimeString);

    // Format date components
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    // Return a concise date/time string
    if (includeDay) {
        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    } else return `${hours}:${minutes}:${seconds}`;
};

interface LogTableProps {
    logsApi: LogsApi;
    selectedLogId: number | null;
    setSelectedLogId?: (id: number | null) => void;
    logFilters: LogFilters;
}

const LogTable: React.FC<LogTableProps> = ({ logsApi: { data, loading, error }, selectedLogId, setSelectedLogId, logFilters }) => {
    const tableStyle: CSSProperties = {};
    if (loading) {
        tableStyle.borderColor = "yellow";
    }
    const toggleRowSelection = (logId: number) => {
        if (setSelectedLogId) {
            if (selectedLogId === logId) {
                setSelectedLogId(null);
            } else {
                setSelectedLogId(logId);
            }
        }
    }
    const getRowStyle = (logId: number): CSSProperties => {
        if (selectedLogId === logId) {
            return {
                backgroundColor: "#34ebc67a"
            }
        }
        return {};
    }
    return (
        <div id="loggerTable" className='flex flex-col gap-2.5 border rounded-md p-1 w-auto'>
            <h2>Logs</h2>
            <div className='overflow-y-auto grow w-auto'>
                <table className="border border-collapse w-200">
                    <thead>
                        <tr>
                            <th className='p-1 border w-30'>Time</th>
                            <th className='p-1 border w-50'>Logger</th>
                            <th className='p-1 border w-20'>level</th>
                            <th className='p-1 border'>Message</th>
                        </tr>
                    </thead>
                    <tbody>
                        {error ? (
                            <tr><td className='p-1 border' colSpan={4}>Error: {error}</td></tr>
                        ) : data && data.length > 0 ? (
                            data?.map((log) => (
                                <tr key={log.id} onClick={() => toggleRowSelection(log.id)} style={getRowStyle(log.id)}>
                                    <td className='p-1 border'>{formatDateString(log.timestamp, false)}</td>
                                    <td className='p-1 border'>{log.logger}</td>
                                    <td className='p-1 border'>{log.level.toUpperCase()}</td>
                                    <td className='p-1 border'>{log.message}</td>
                                </tr>
                            ))
                        ) : (
                            <tr><td className='p-1 border' colSpan={4}>No logs to display</td></tr>
                        )}
                    </tbody>
                </table>
            </div>
            <Pagination logFilters={logFilters} />
        </div>
    )
}

export default LogTable;
