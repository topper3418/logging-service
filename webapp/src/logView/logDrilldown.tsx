import React from 'react';
import { useFetchLogData } from './logHooks';


interface LogDrilldownProps {
    logId: number | null;
}

const LogDrilldown: React.FC<LogDrilldownProps> = ({ logId }) => {
    const { data, loading, error } = useFetchLogData(logId);
    return (
        <div id="logDrilldown" className='flex flex-col gap-2.5 border rounded-md p-1 grow'>
            <h2>Log Drilldown</h2>
            {!logId ? <p>Select Log...</p> :
                loading ? <p>Loading...</p> :
                    error ? <p>Error: {error}</p> :
                        data && <>
                            <p>Log ID: {logId}</p>
                            <p>Timestamp: {data?.timestamp}</p>
                            <p>Logger: {data?.logger}</p>
                            <p>Level: {data?.level}</p>
                            <p>Message: {data?.message}</p>
                            <p>Meta: </p>
                            <div className='grow overflow-y-auto'>
                                <pre>{JSON.stringify(data?.meta, null, 2)}</pre>
                            </div>
                        </>}
        </div>
    )
}



export default LogDrilldown;


