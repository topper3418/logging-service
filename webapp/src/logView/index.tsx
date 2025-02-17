import React, { useEffect, useState } from "react"
import { useFetchLogs, useFilters } from "./logHooks";
import Banner from "./banner";
import Filters from "./filters";
import LogTable from "./logTable";
import { useFetchLoggers } from "./loggerHooks";
import LoggerTable from "./loggerTable";
import LogDrilldown from "./logDrilldown";


const LogView: React.FC = () => {
    const [selectedLogId, setSelectedLogId] = useState<number | null>(null);
    const [showLoggers, setShowLoggers] = useState<boolean>(false);
    const logFilters = useFilters();

    const logs = useFetchLogs(logFilters.get);
    const loggers = useFetchLoggers();

    const refreshView = () => {
        logs.refetch();
        loggers.refetch();
    }
    useEffect(() => {
        refreshView()
    }, [
        logFilters.get.offset,
        logFilters.get.limit,
        logFilters.get.excludeLoggers,
        logFilters.get.search,
        logFilters.get.minTime,
        logFilters.get.maxTime,
    ])
    return (
        <div id="logView" className="flex h-full flex-col p-2.5 gap-2.5 max-h-full">
            <Banner title="Log View">
                <></>
                <button onClick={refreshView}>Refresh</button>
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
