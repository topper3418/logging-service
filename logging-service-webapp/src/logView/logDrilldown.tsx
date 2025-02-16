import React from 'react';
import { useFetchLogData } from './logHooks';


interface LogDrilldownProps {
    logId: number | null;
}

const logDrilldownStyle: React.CSSProperties = {
    overflowY: "auto",
    flexGrow: 1,
}

const LogDrilldown: React.FC<LogDrilldownProps> = ({ logId }) => {
    const { data, loading, error } = useFetchLogData(logId);
    return (
        <div id="logDrilldown" style={logDrilldownStyle}>
            <h1>Log Drilldown</h1>
            {loading ? <p>Loading...</p> :
                error ? <p>Error: {error}</p> :
                    data && <>
                        <p>Log ID: {logId}</p>
                        <p>Timestamp: {data?.timestamp}</p>
                        <p>Logger: {data?.logger}</p>
                        <p>Level: {data?.level}</p>
                        <p>Message: {data?.message}</p>
                        <p>Meta: </p>
                        <div>
                            <pre>{JSON.stringify(data?.meta, null, 2)}</pre>
                        </div>
                    </>}
        </div>
    )
}



export default LogDrilldown;


