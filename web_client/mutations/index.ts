import axios from "axios"

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export const handleCreateOrganization = async(organization: any) => {
  try {
    const response = await axios.post(`${API_URL}/organization/create/`, organization)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
};

export const handleCreateFolder = async (folder: any) => {
  try {
    const response = await axios.post(`${API_URL}/folder/create/`, folder);
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
};

export const handleCreateFile = async (file: any) => {
  try {
    const response = await axios.post(`${API_URL}/file/create/`, file);
    return response
  } catch (error: any) {
    console.error("Error in handleCreateFile:", error.response ? error.response.data : error.message);
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
};

export const handleUploadFile = async (file: any) => {
  try {
    const response = await axios.post(`/api/upload/file/`, file);
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
};

export const handleRenameFolder = async (folder: any, folderId: string) => {
  try {
    console.log("Function called: ", folderId)
    const response = await axios.put(`${API_URL}/folder/update/${folderId}`, folder)
    console.log("Response: ", response)
    return response
  } catch (error: any) {
    console.error('Error updating folder:', error);
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
}

export const handleRenameFile = async (file: any, fileId: string) => {
  try {
    const response = await axios.put(`${API_URL}/file/update/${fileId}`, file)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
}

export const handleRenameFleet = async (fleet: any, fleetId: string) => {
  try {
    const response = await axios.put(`${API_URL}/fleet/update/${fleetId}`, fleet)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
}

export const handleMoveFolderToTrash = async (folderId: string) => {
  try {
    const response = await axios.put(`${API_URL}/folder/delete/${folderId}`)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleMoveFileToTrash = async (fileId: string) => {
  try {
    const response = await axios.put(`${API_URL}/file/delete/${fileId}`)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleDeleteFolder = async (folderId: string) => {
  try {
    const response = await axios.delete(`${API_URL}/folder/delete/permanent/${folderId}`)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleDeleteFile = async (fileId: string) => {
  try {
    const response = await axios.delete(`${API_URL}/file/delete/permanent/${fileId}`)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
}

export const handleRestoreFolder = async (folderId: string) => {
  try {
    const response = await axios.put(`${API_URL}/folder/restore/${folderId}`)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleRestoreFile = async (fileId: string) => {
  try {
    const response = await axios.put(`${API_URL}/file/restore/${fileId}`)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleUpdateUser = async (email: string | undefined, user: any) => {
  try {
    const response = await axios.put(`${API_URL}/auth/update/${email}`, user)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleCreateDevice = async (device: any) => {
  try {
    const response = await axios.post(`${API_URL}/device/create/`, device)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleCreateFleet = async (fleet: any) => {
  try {
    const response = await axios.post(`${API_URL}/fleet/create/`, fleet)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
}

export const handleUpdateDevice = async (deviceId: string ,device: any) => {
  try {
    const response = await axios.put(`${API_URL}/device/update/${deviceId}`, device)
    return response
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } };
  }
}

export const handleEmptyBin = async (organizationId: string) => {
  try {
    if (organizationId) {
      const response = await axios.delete(`${API_URL}/bin/empty/${organizationId}`)
      return response  
    }
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleDeleteFleet = async (fleetId: string) => {
  try {
    if (fleetId) {
      const response = await axios.delete(`${API_URL}/fleet/delete/${fleetId}`)
      return response  
    }
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}

export const handleDeleteDevice = async (deviceId: string) => {
  try {
    if (deviceId) {
      const response = await axios.delete(`${API_URL}/device/delete/${deviceId}`)
      return response  
    }
  } catch (error: any) {
    return error.response ? error.response : { data: { error: "Unknown error occurred" } }; 
  }
}
