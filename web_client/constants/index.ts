"use client";

import { useParams } from "next/navigation";
import { saveAs } from "file-saver";
import JSZip from 'jszip';
import { fetchFolderHierarchyRecursively } from "@/queries";
import axios from "axios";

export const API_URL = process.env.NEXT_PUBLIC_API_URL;

export const useOrganizationId = () => {
  const params = useParams();
  const { organization_id } = params;
  const orgId = Array.isArray(organization_id) ? organization_id[0] : organization_id;
  return orgId;
};

export const useFolderId = () => {
  const params = useParams();
  const { folder_id } = params;
  const folderId = Array.isArray(folder_id) ? folder_id[0] : folder_id ?? null;
  return folderId;
};

export const transformHierarchy = (folder: any) => {
  return {
    name: folder.name,
    created_at: folder.created_at,
    updated_at: folder.updated_at,
    files: folder.files?.map((file: any) => ({
      name: file.name,
      file_path: file.file_path,
      file_size: file.file_size,
      created_at: file.created_at,
      updated_at: file.updated_at,
    })),
    subfolders: folder.subfolders?.map((subfolder: any) => transformHierarchy(subfolder)),
  };
};

export const formatFileSize = (sizeInBytes: number): string => {
  if (sizeInBytes < 1024) {
    return `${sizeInBytes} bytes`;
  } else if (sizeInBytes < 1024 * 1024) {
    return `${(sizeInBytes / 1024).toFixed(2)} KB`;
  } else if (sizeInBytes < 1024 * 1024 * 1024) {
    return `${(sizeInBytes / (1024 * 1024)).toFixed(2)} MB`;
  } else {
    return `${(sizeInBytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
  }
};

export const generateManifestForDownload = (hierarchy: any, folderName: string | null = null) => {
  const transformedHierarchy = transformHierarchy(hierarchy);
  console.log("Transformed Hierachy: ", transformedHierarchy)
  const manifest = JSON.stringify(transformedHierarchy, null, 2);
  const fileName = folderName && folderName !== "root" ? (`${folderName} manifest.json`) : ('drive manifest.json');
  const blob = new Blob([manifest], { type: 'application/json' });
  const url = URL.createObjectURL(blob);

  // Create a link element for downloading
  const a = document.createElement('a');
  a.href = url;
  a.download = fileName;
  a.click();

  // Cleanup
  URL.revokeObjectURL(url);
};

export const downloadFolderAsZip = async (organizationId: string, folderId: string = 'root') => {
  try {
    const zip = new JSZip();
    const folderHierarchy = await fetchFolderHierarchyRecursively(organizationId, folderId);

    const addFolderToZip = async (zipFolder: any, folderData: any) => {
      // Add files to the current folder
      for (const file of folderData.files || []) { // Safeguard for empty or undefined files array
        const fileData = await fetchFileData(file.file_path);
        zipFolder.file(file.name, fileData);
      }
    
      // Recursively add subfolders
      if (folderData.subfolders && folderData.subfolders.length > 0) { // Safeguard for undefined or empty subfolders
        for (const subfolder of folderData.subfolders) {
          const newZipFolder = zipFolder.folder(subfolder.name);
          await addFolderToZip(newZipFolder, subfolder);
        }
      }
    };
    
    // Create root folder in the zip
    const rootZipFolder = zip.folder(folderHierarchy.name);
    await addFolderToZip(rootZipFolder, folderHierarchy);

    // Generate zip file and trigger download
    const content = await zip.generateAsync({ type: 'blob' });
    saveAs(content, `${folderHierarchy.name}.zip`);

  } catch (error) {
    console.error("Error downloading folder as zip: ", error);
  }
};

const fetchFileData = async (fileUrl: string) => {
  try {
    const response = await axios.get(fileUrl, {
      responseType: 'blob',
    });
    return response.data; // Return the blob directly
  } catch (error) {
    console.error('Error fetching the file:', error);
    throw error;
  }
};

export const downloadFile = async (fileUrl: string, fileName: string) => {
  try {
    const response = await axios.get(fileUrl, {
      responseType: 'blob',
    });

    const url = window.URL.createObjectURL(response.data);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', fileName);
    document.body.appendChild(link);
    link.click();

    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('Error downloading the file:', error);
  }
};
