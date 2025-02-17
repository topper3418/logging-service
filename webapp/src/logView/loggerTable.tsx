import React, { CSSProperties } from 'react';
import { Logger, LoggersApi, useSetLoggerLevel } from './loggerHooks';
import { LogFilters } from './logHooks';


interface LoggerTableProps {
    loggersApi: LoggersApi;
    logFilters: LogFilters;
    refreshView: () => void;
    hideLoggers: () => void;
}

const LoggerTable: React.FC<LoggerTableProps> = ({ loggersApi: { data, loading, error }, logFilters, refreshView, hideLoggers }) => {
    // value being true means "include", false means "exclude"
    const setLoggerExclusion = (excludeLogger: Logger, value: boolean) => {
        console.log(`setting logger filter: ${excludeLogger.name}, value: ${value}`);
        if (value) {
            logFilters.set.excludeLoggers.remove(excludeLogger.id);
        } else {
            logFilters.set.excludeLoggers.add(excludeLogger.id);
        }
        refreshView();
    }
    const allChecked = logFilters.get.excludeLoggers.length === 0;
    const allLoggerIds = data?.map((logger) => logger.id);
    const handleMasterCheckbox = () => {
        if (allChecked) {
            console.log("all checked, unchecking all");
            logFilters.set.excludeLoggers.raw(allLoggerIds);
        } else {
            console.log("not all checked, checking all");
            logFilters.set.excludeLoggers.raw([]);
        }
        refreshView();
    }
    const tableStyle: CSSProperties = {};
    if (loading) {
        tableStyle.borderColor = "yellow";
    }
    const sortedData = data?.sort((a, b) => a.name.localeCompare(b.name));
    return (
        <div id="loggerTable" className='flex flex-col gap-2.5 w-auto max-w-1/3 border rounded-md p-1 '>
            <div className='flex justify-between'>
                <h2>Loggers</h2>
                <button onClick={hideLoggers}>Hide</button>
            </div>
            <div className='overflow-y-auto'>
                <table className='border border-collapse'>
                    <thead>
                        <tr>
                            <th className='p-1 border'>
                                <input
                                    type="checkbox"
                                    checked={allChecked}
                                    onChange={handleMasterCheckbox} />
                            </th>
                            <th className='p-1 border'>Logger</th>
                            <th className='p-1 border'>Level</th>
                        </tr>
                    </thead>
                    <tbody>
                        {error ? (
                            <tr><td colSpan={4} className='p-1 border'>Error: {error}</td></tr>
                        ) : sortedData && sortedData.length > 0 ? sortedData?.map((logger) => (
                            <tr key={logger.id}>
                                <td className='p-1 border'>
                                    <input
                                        type="checkbox"
                                        checked={!logFilters.get.excludeLoggers.includes(logger.id)}
                                        onChange={(event) => setLoggerExclusion(logger, event.target.checked)} />
                                </td>
                                <td className='p-1 border'>{logger.name}</td>
                                <td className='p-1 border'>
                                    <LoggerLevelDropdown
                                        loggerId={logger.id}
                                        currentLevel={logger.level}
                                        refreshCallback={refreshView} />
                                </td>
                            </tr>
                        )) : (
                            <tr><td colSpan={4} className='p-1 border'>No loggers to display</td></tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    )
}

interface LoggerLevelDropdownProps {
    loggerId: number;
    currentLevel: string;
    refreshCallback: () => void;
}

const LoggerLevelDropdown: React.FC<LoggerLevelDropdownProps> = ({ loggerId, currentLevel, refreshCallback }) => {
    const [setLoggerLevel, { loading, error }] = useSetLoggerLevel();
    const levels = ['debug', 'info', 'warn', 'error'];
    const changeLevel = (level: string) => {
        console.log(`changing level for logger id ${loggerId} to ${level}`);
        setLoggerLevel(loggerId, level, refreshCallback);
    }
    return (<>
        {loading ? <p>Setting level...</p> :
            error ? <p>Error: {error}</p> :
                <select value={currentLevel} onChange={(e) => changeLevel(e.target.value)}>
                    {levels.map((level) => (
                        <option key={level} value={level}>{level.toUpperCase()}</option>
                    ))}
                </select>}
    </>)
}

export default LoggerTable;
