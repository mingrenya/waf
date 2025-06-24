import React, { useState, useEffect } from 'react';

interface ThreatData {
  id: string;
  severity: 'low' | 'medium' | 'high';
  description: string;
  modelAnalysis: string;
  timestamp: Date;
}

const ThreatAnalysis: React.FC = () => {
  const [threats, setThreats] = useState<ThreatData[]>([]);

  useEffect(() => {
    // 从后端API获取威胁数据
    const fetchThreats = async () => {
      const response = await fetch('/api/threats');
      const data = await response.json();
      setThreats(data);
    };
    
    fetchThreats();
    const interval = setInterval(fetchThreats, 30000); // 每30秒更新
    
    return () => clearInterval(interval);
  }, []);

  const getSeverityColor = (severity: string) => {
    switch(severity) {
      case 'high': return 'bg-red-500';
      case 'medium': return 'bg-yellow-500';
      default: return 'bg-green-500';
    }
  };

  return (
    <div className="p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-xl font-bold mb-4 text-security-blue">实时威胁分析</h2>
      <div className="space-y-4">
        {threats.map(threat => (
          <div key={threat.id} className="border-l-4 border-red-500 pl-4 py-2">
            <div className="flex items-center justify-between">
              <span className={`inline-block w-3 h-3 rounded-full mr-2 ${getSeverityColor(threat.severity)}`}></span>
              <span className="font-medium">{threat.description}</span>
              <span className="text-sm text-gray-500">
                {new Date(threat.timestamp).toLocaleTimeString()}
              </span>
            </div>
            <div className="mt-2 text-sm bg-gray-50 p-3 rounded">
              <span className="font-semibold">模型分析:</span> 
              <p className="mt-1">{threat.modelAnalysis}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ThreatAnalysis;
