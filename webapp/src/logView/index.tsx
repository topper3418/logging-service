import React, { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom";
import { useFetchLogs, useFilters } from "./logHooks";
import Banner from "./banner";
import Filters from "./filters";
import LogTable from "./logTable";
import { useFetchLoggers } from "./loggerHooks";
import LoggerTable from "./loggerTable";
import LogDrilldown from "./logDrilldown";


const LogView: React.FC = () => {
    const navigate = useNavigate();
    const [selectedLogId, setSelectedLogId] = useState<number | null>(null);
    const [showLoggers, setShowLoggers] = useState<boolean>(true);
    const [autoRefetchEnabled, setAutoRefetchEnabled] = useState<boolean>(true);
    const logFilters = useFilters();

    const logs = useFetchLogs(logFilters.get);
    const loggers = useFetchLoggers();

    const refreshView = () => {
        logs.refetch();
        loggers.refetch();
    }
    const autoRefetch = () => {
        if (!autoRefetchEnabled) return
        refreshView()
    }
    useEffect(() => {
        const timeout = setTimeout(autoRefetch, 500)
        return () => clearTimeout(timeout)
    }, [logs.loading])
    useEffect(() => {
        refreshView()
    }, [logFilters])
    return (
        <div id="logView" className="flex h-full flex-col p-2.5 gap-2.5 max-h-full">
            <Banner title="Log View">
                <button onClick={() => navigate("/")}>Back</button>
                <button onClick={() => setAutoRefetchEnabled(!autoRefetchEnabled)}>
                    {autoRefetchEnabled ? 'Refresh Enabled' : "Refresh Disabled"}
                </button>
            </Banner>
            <Filters logFilters={logFilters} />
            <div className="flex flex-row gap-2.5 m-0 grow overflow-hidden justify-start">
                {showLoggers ?
                    <LoggerTable
                        loggersApi={loggers}
                        logFilters={logFilters}
                        refreshView={refreshView}
                        hideLoggers={() => setShowLoggers(false)} /> :
                    <div>
                        <button
                            onClick={() => setShowLoggers(true)}>
                            Loggers &#x25B6;
                        </button>
                    </div>}
                <LogTable
                    logsApi={logs}
                    selectedLogId={selectedLogId}
                    setSelectedLogId={setSelectedLogId}
                    logFilters={logFilters} />
                <LogDrilldown logId={selectedLogId} />
            </div>
        </div >
    )
}


export default LogView;
